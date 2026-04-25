#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>
#include <TinyGPS++.h>
#include "MAX30105.h"        
#include "heartRate.h"       
#include <Adafruit_MLX90614.h>

// --- IDENTITAS ALAT PENDAKI ---
String nodeID = "PENDAKI_01"; 

// --- KONFIGURASI PIN ESP32-C3 ---
#define I2C_SDA 10
#define I2C_SCL 2
#define GPS_RX 4 
#define GPS_TX 5
#define LORA_RX 6 
#define LORA_TX 7

Adafruit_SSD1306 display(128, 32, &Wire, -1);
TinyGPSPlus gps;
HardwareSerial gpsSerial(1);  // UART1 untuk GPS
HardwareSerial loraSerial(0); // UART0 untuk LoRa
MAX30105 particleSensor;
Adafruit_MLX90614 mlx = Adafruit_MLX90614();

String strBPM = "--", strTemp = "--", strLat = "--", strLng = "--";
const byte RATE_SIZE = 4;
byte rates[RATE_SIZE]; 
byte rateSpot = 0;
long lastBeat = 0; 
float beatsPerMinute;
int beatAvg = 0;

unsigned long lastSendTime = 0;
unsigned long sendInterval = 3000;

void setup() {
  Serial.begin(115200);
  
  // Memulai jalur I2C untuk 3 perangkat sekaligus (OLED, MLX, MAX)
  Wire.begin(I2C_SDA, I2C_SCL);

  gpsSerial.begin(9600, SERIAL_8N1, GPS_RX, GPS_TX);
  loraSerial.begin(9600, SERIAL_8N1, LORA_RX, LORA_TX);

  // Inisialisasi OLED
  if(!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) {
    Serial.println(F("OLED Gagal!"));
    while(1); 
  }
  
  display.clearDisplay(); display.setTextSize(1); display.setTextColor(SSD1306_WHITE);
  display.setCursor(15, 12); display.println(F("MEMULAI SENSOR..."));
  display.display();
  
  // Inisialisasi Sensor
  particleSensor.begin(Wire, I2C_SPEED_STANDARD); // Ingat, harus Standard agar MLX tidak NaN
  particleSensor.setup(); 
  particleSensor.setPulseAmplitudeRed(0x0A); 
  mlx.begin();
  
  // Seed untuk jeda waktu acak anti-tabrakan sinyal LoRa
  randomSeed(analogRead(0)); 
}

void loop() {
  // 1. BACA GPS
  while (gpsSerial.available() > 0) {
    gps.encode(gpsSerial.read());
  }

  // 2. BACA DETAK JANTUNG
  long irValue = particleSensor.getIR();
  if (checkForBeat(irValue) == true) {
    long delta = millis() - lastBeat;
    lastBeat = millis();
    beatsPerMinute = 60 / (delta / 1000.0);
    if (beatsPerMinute < 255 && beatsPerMinute > 20) {
      rates[rateSpot++] = (byte)beatsPerMinute; 
      rateSpot %= RATE_SIZE;
      beatAvg = 0;
      for (byte x = 0 ; x < RATE_SIZE ; x++) beatAvg += rates[x];
      beatAvg /= RATE_SIZE;
    }
  }

  // 3. PENGIRIMAN & UPDATE LAYAR OLED
  if (millis() - lastSendTime > sendInterval) {
    lastSendTime = millis();
    sendInterval = random(3000, 5000); // Acak waktu kirim antara 3-5 detik
    
    double suhuTubuh = mlx.readObjectTempC();
    strBPM = (irValue < 50000) ? "--" : String(beatAvg);
    strTemp = String(suhuTubuh, 1);
    
    if (gps.location.isValid()) {
      strLat = String(gps.location.lat(), 6);
      strLng = String(gps.location.lng(), 6);
    } else {
      strLat = "Mencari.."; 
      strLng = "Mencari..";
    }

    // Tembak via LoRa
    String payload = nodeID + "|" + strBPM + "|" + strTemp + "|" + strLat + "|" + strLng;
    loraSerial.println(payload);
    
    // Tampilkan ke layar OLED Pendaki
    display.clearDisplay(); 
    display.setTextSize(1);
    
    display.setCursor(0, 0); 
    display.print("TX | "); display.print(nodeID);
    
    display.setCursor(0, 8); 
    display.print("HR:"); display.print(strBPM); display.print(" T:"); display.print(strTemp);
    
    display.setCursor(0, 16); 
    display.print("Lat: "); display.print(strLat);
    
    display.setCursor(0, 24); 
    display.print("Lng: "); display.print(strLng);
    
    display.display();
  }
}