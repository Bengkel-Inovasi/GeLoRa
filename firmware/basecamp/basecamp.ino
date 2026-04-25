/*
 * GeLoRa Basecamp Receiver Firmware
 * Hardware: ESP32 (any variant) + LoRa SX1276/SX1278 + WiFi
 *
 * Role:
 *   1. Listens on LoRa for wristband JSON packets
 *   2. Forwards them to the MQTT broker on topic /server/record
 *   3. Stays permanently installed at the mountain basecamp
 *
 * Libraries required (install via Arduino Library Manager):
 *   - LoRa (by Sandeep Mistry)
 *   - PubSubClient (by Nick O'Leary)
 *   - ArduinoJson
 *
 * Wiring (standard ESP32 DevKit, adjust pins for your board):
 *   LoRa SX1276:
 *     NSS/CS  → GPIO5
 *     RESET   → GPIO14
 *     DIO0    → GPIO2
 *     MOSI    → GPIO23
 *     MISO    → GPIO19
 *     SCK     → GPIO18
 */

#include <SPI.h>
#include <LoRa.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>

// ─── Network configuration ──────────────────────────────────────────────────
#define WIFI_SSID     "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"

// ─── MQTT broker ─────────────────────────────────────────────────────────────
// IP of the machine running Mosquitto (same machine as the GeLoRa backend)
#define MQTT_HOST     "192.168.1.100"
#define MQTT_PORT     1883
#define MQTT_USER     "gelora-mqtt"
#define MQTT_PASSWORD "gelora-mqtt"
#define MQTT_CLIENT   "gelora-basecamp"
#define MQTT_TOPIC    "/server/record"

// ─── LoRa pins ───────────────────────────────────────────────────────────────
#define LORA_NSS    5
#define LORA_RST    14
#define LORA_DIO0   2
#define LORA_FREQ   915E6  // must match wristband

// ─── Globals ─────────────────────────────────────────────────────────────────
WiFiClient   wifiClient;
PubSubClient mqtt(wifiClient);

// ─── WiFi connection ─────────────────────────────────────────────────────────
void connectWiFi() {
  Serial.print("Connecting to WiFi");
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println();
  Serial.print("WiFi connected, IP: ");
  Serial.println(WiFi.localIP());
}

// ─── MQTT connection ──────────────────────────────────────────────────────────
void connectMQTT() {
  mqtt.setServer(MQTT_HOST, MQTT_PORT);
  while (!mqtt.connected()) {
    Serial.print("Connecting to MQTT broker...");
    if (mqtt.connect(MQTT_CLIENT, MQTT_USER, MQTT_PASSWORD)) {
      Serial.println(" connected");
    } else {
      Serial.print(" failed (rc=");
      Serial.print(mqtt.state());
      Serial.println("), retry in 5s");
      delay(5000);
    }
  }
}

// ─── Setup ───────────────────────────────────────────────────────────────────
void setup() {
  Serial.begin(115200);

  connectWiFi();
  connectMQTT();

  // LoRa — must use same frequency and modulation as the wristband
  LoRa.setPins(LORA_NSS, LORA_RST, LORA_DIO0);
  if (!LoRa.begin(LORA_FREQ)) {
    Serial.println("LoRa init failed");
    while (true) delay(10);
  }
  LoRa.setSpreadingFactor(9);
  LoRa.setSignalBandwidth(125E3);
  LoRa.setCodingRate4(5);
  Serial.println("GeLoRa basecamp receiver ready");
}

// ─── Main loop ────────────────────────────────────────────────────────────────
void loop() {
  // Keep MQTT alive
  if (!mqtt.connected()) connectMQTT();
  mqtt.loop();

  // Keep WiFi alive
  if (WiFi.status() != WL_CONNECTED) connectWiFi();

  int packetSize = LoRa.parsePacket();
  if (packetSize == 0) return;

  // Read LoRa packet
  char buf[256] = {0};
  int i = 0;
  while (LoRa.available() && i < (int)sizeof(buf) - 1) {
    buf[i++] = (char)LoRa.read();
  }
  buf[i] = '\0';

  Serial.print("LoRa received (RSSI=");
  Serial.print(LoRa.packetRssi());
  Serial.print("dBm): ");
  Serial.println(buf);

  // Validate: must be a JSON object with a "mid" field
  StaticJsonDocument<256> doc;
  DeserializationError err = deserializeJson(doc, buf);
  if (err || !doc.containsKey("mid")) {
    Serial.println("Ignored: not valid GeLoRa payload");
    return;
  }

  // Forward to MQTT broker as-is (backend expects the same JSON)
  bool ok = mqtt.publish(MQTT_TOPIC, buf);
  Serial.println(ok ? "Published to MQTT" : "MQTT publish failed");
}
