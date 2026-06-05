/*
 * GeLoRa Wristband Firmware
 * Hardware: ESP32-C3 + MAX30102 + MLX90614 + BN220 GPS + LoRa SX1276/SX1278
 *
 * Libraries required (install via Arduino Library Manager):
 *   - SparkFun MAX3010x Pulse and Proximity Sensor Library
 *   - Adafruit MLX90614 Library
 *   - TinyGPSPlus
 *   - LoRa (by Sandeep Mistry)
 *   - ArduinoJson
 *
 * Wiring (ESP32-C3 default pins, adjust as needed):
 *   MAX30102  → SDA=GPIO8, SCL=GPIO9  (I2C)
 *   MLX90614  → SDA=GPIO8, SCL=GPIO9  (I2C, shared bus)
 *   BN220 GPS → RX=GPIO20, TX=GPIO21  (UART1)
 *   LoRa SX1276:
 *     NSS/CS  → GPIO7
 *     RESET   → GPIO3
 *     DIO0    → GPIO2
 *     MOSI    → GPIO6
 *     MISO    → GPIO5
 *     SCK     → GPIO4
 */

#include <Wire.h>
#include <SPI.h>
#include <LoRa.h>
#include <ArduinoJson.h>

// MAX30102
#include "MAX30105.h"
#include "heartRate.h"

// MLX90614
#include <Adafruit_MLX90614.h>

// GPS
#include <TinyGPSPlus.h>
#include <HardwareSerial.h>

// ─── Device identity ────────────────────────────────────────────────────────
// This MID must match the Node registered/validated in the GeLoRa dashboard.
// Use a unique string per device (e.g. based on MAC address).
#define DEVICE_MID "gelora-wristband-001"

// ─── LoRa pins ──────────────────────────────────────────────────────────────
#define LORA_NSS    7
#define LORA_RST    3
#define LORA_DIO0   2
#define LORA_FREQ   915E6   // 915 MHz for SEA/US; use 868E6 for EU

// ─── GPS UART ────────────────────────────────────────────────────────────────
#define GPS_RX 20
#define GPS_TX 21
#define GPS_BAUD 9600

// ─── Timing ─────────────────────────────────────────────────────────────────
#define SEND_INTERVAL_MS 5000   // send reading every 5 seconds

// ─── Filtering Constants ────────────────────────────────────────────────────
#define TEMP_FILTER_SIZE 10
#define BPM_FILTER_SIZE  8

// ─── Globals ────────────────────────────────────────────────────────────────
MAX30105         maxSensor;
Adafruit_MLX90614 mlx;
TinyGPSPlus      gps;
HardwareSerial   gpsSerial(1);

// Heart rate algorithm state
const byte  HR_SAMPLE_SIZE = 4;
byte        hrRates[HR_SAMPLE_SIZE];
byte        hrRateSpot = 0;
long        lastBeatTime = 0;
float       beatsPerMinute = 0;
float       beatAvg = 0;

// Moving Average Filters
float tempHistory[TEMP_FILTER_SIZE];
int tempIndex = 0;
float bpmHistory[BPM_FILTER_SIZE];
int bpmIndex = 0;

unsigned long lastSendMs = 0;

// ─── Filter Functions ───────────────────────────────────────────────────────
float movingAverage(float newValue, float *history, int *index, int size) {
  history[*index] = newValue;
  *index = (*index + 1) % size;
  
  float sum = 0;
  int count = 0;
  for (int i = 0; i < size; i++) {
    if (history[i] > 0) {
      sum += history[i];
      count++;
    }
  }
  return (count > 0) ? (sum / count) : 0;
}

// ─── Setup ──────────────────────────────────────────────────────────────────
void setup() {
  Serial.begin(115200);
  Wire.begin(8, 9);  // SDA, SCL

  // Initialize filter histories
  for (int i = 0; i < TEMP_FILTER_SIZE; i++) tempHistory[i] = 0;
  for (int i = 0; i < BPM_FILTER_SIZE; i++) bpmHistory[i] = 0;

  // MAX30102
  if (!maxSensor.begin(Wire, I2C_SPEED_FAST)) {
    Serial.println("MAX30102 not found");
    while (true) delay(10);
  }
  maxSensor.setup();
  // Increased amplitude for better penetration on wrist-top
  maxSensor.setPulseAmplitudeRed(0x1F); 
  maxSensor.setPulseAmplitudeIR(0x1F);
  maxSensor.setPulseAmplitudeGreen(0);
  Serial.println("MAX30102 ready");

  // MLX90614
  if (!mlx.begin()) {
    Serial.println("MLX90614 not found");
    while (true) delay(10);
  }
  Serial.println("MLX90614 ready");

  // BN220 GPS
  gpsSerial.begin(GPS_BAUD, SERIAL_8N1, GPS_RX, GPS_TX);
  Serial.println("GPS UART started");

  // LoRa
  LoRa.setPins(LORA_NSS, LORA_RST, LORA_DIO0);
  if (!LoRa.begin(LORA_FREQ)) {
    Serial.println("LoRa init failed");
    while (true) delay(10);
  }
  LoRa.setSpreadingFactor(9);
  LoRa.setSignalBandwidth(125E3);
  LoRa.setCodingRate4(5);
  LoRa.setTxPower(17);
  Serial.println("LoRa ready");

  Serial.println("GeLoRa wristband ready — MID: " DEVICE_MID);
}

// ─── Heart rate sampling ─────────────────────────────────────────────────────
void sampleHeartRate() {
  long irValue = maxSensor.getIR();
  if (checkForBeat(irValue)) {
    long delta = millis() - lastBeatTime;
    lastBeatTime = millis();
    beatsPerMinute = 60.0 / (delta / 1000.0);
    if (beatsPerMinute > 20 && beatsPerMinute < 255) {
      hrRates[hrRateSpot++] = (byte)beatsPerMinute;
      hrRateSpot %= HR_SAMPLE_SIZE;
      beatAvg = 0;
      for (byte i = 0; i < HR_SAMPLE_SIZE; i++) beatAvg += hrRates[i];
      beatAvg /= HR_SAMPLE_SIZE;
    }
  }
}

// ─── GPS feeding ─────────────────────────────────────────────────────────────
void feedGPS() {
  while (gpsSerial.available()) gps.encode(gpsSerial.read());
}

// ─── Main loop ───────────────────────────────────────────────────────────────
void loop() {
  sampleHeartRate();
  feedGPS();

  if (millis() - lastSendMs < SEND_INTERVAL_MS) return;
  lastSendMs = millis();

  // Build JSON payload
  StaticJsonDocument<256> doc;
  doc["mid"] = DEVICE_MID;

  // Heart rate with Moving Average Filter
  if (beatAvg > 20) {
    float filteredBpm = movingAverage(beatAvg, bpmHistory, &bpmIndex, BPM_FILTER_SIZE);
    doc["heart_rate"] = round(filteredBpm * 10) / 10.0;
  }

  // Body temperature with Moving Average Filter
  float bodyTemp = mlx.readObjectTempC();
  if (!isnan(bodyTemp) && bodyTemp > 30.0 && bodyTemp < 42.0) {
    float filteredTemp = movingAverage(bodyTemp, tempHistory, &tempIndex, TEMP_FILTER_SIZE);
    doc["body_temperature"] = round(filteredTemp * 10) / 10.0;
  }

  // GPS
  if (gps.location.isValid() && gps.location.age() < 2000) {
    doc["latitude"]  = gps.location.lat();
    doc["longitude"] = gps.location.lng();
  }

  char payload[256];
  serializeJson(doc, payload);

  // Send via LoRa
  LoRa.beginPacket();
  LoRa.print(payload);
  LoRa.endPacket();

  Serial.print("Sent: ");
  Serial.println(payload);
}
