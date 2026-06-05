#include <WiFi.h>
#include <PubSubClient.h>
#include <WiFiManager.h> // Required: Install "WiFiManager" by tzapu

// --- KONFIGURASI MQTT SERVER (UBUNTU) ---
const char* mqtt_server = "10.127.41.105"; // WSL2 IP (mirrored mode)
const int mqtt_port = 1883;
const char* mqtt_user = "gelora-mqtt";  // USERNAME YANG ANDA BUAT TADI
const char* mqtt_pass = "gelora-mqtt";   // PASSWORD YANG ANDA BUAT TADI
const char* mqtt_topic = "/server/record";

// --- KONFIGURASI PIN LORA ---
#define LORA_RX_PIN 26 
#define LORA_TX_PIN 27 
HardwareSerial loraSerial(2);

WiFiClient espClient;
PubSubClient client(espClient);
WiFiManager wm;

void setup_wifi() {
  // Starts an Access Point "GeLoRa-Basecamp-Example" if it cannot connect to saved WiFi
  bool res = wm.autoConnect("GeLoRa-Basecamp-Example");

  if(!res) {
    Serial.println("Gagal konek ke WiFi!");
    ESP.restart();
  } else {
    Serial.println("\nWiFi Terhubung! IP: ");
    Serial.println(WiFi.localIP());
  }
}

void reconnect() {
  // Looping sampai berhasil konek ke MQTT Ubuntu
  while (!client.connected()) {
    Serial.print("Mencoba konek ke MQTT Ubuntu...");
    String clientId = "ESP32_Basecamp_Gateway";
    
    // DISINILAH KITA MEMASUKKAN USERNAME & PASSWORD
    if (client.connect(clientId.c_str(), mqtt_user, mqtt_pass)) {
      Serial.println("BERHASIL!");
    } else {
      Serial.print("Gagal, error kode=");
      Serial.print(client.state());
      Serial.println(" Coba lagi dalam 5 detik");
      delay(5000);
    }
  }
}

void setup() {
  Serial.begin(115200);
  loraSerial.begin(9600, SERIAL_8N1, LORA_RX_PIN, LORA_TX_PIN);

  setup_wifi();
  client.setServer(mqtt_server, mqtt_port);
  Serial.println(F("=== GATEWAY RX SIAP ==="));
}

String toJsonNum(String val) {
  if (val.length() == 0) return "null";
  char c = val.charAt(0);
  if ((c >= '0' && c <= '9') || c == '-') return val;
  return "null";
}

void loop() {
  if (!client.connected()) reconnect();
  client.loop();

  // JIKA ADA DATA LORA MASUK DARI PENDAKI
  if (loraSerial.available() > 0) {
    String incoming = loraSerial.readStringUntil('\n');
    incoming.trim();

    if (incoming.length() > 0) {
      Serial.println("LoRa Masuk: " + incoming);

      // Parse Data: PENDAKI_01|78|36.2|-7.79|110.36
      int p1 = incoming.indexOf('|');
      int p2 = incoming.indexOf('|', p1 + 1);
      int p3 = incoming.indexOf('|', p2 + 1);
      int p4 = incoming.indexOf('|', p3 + 1);

      if (p1 != -1 && p4 != -1) {
        String id   = incoming.substring(0, p1);
        String bpm  = incoming.substring(p1 + 1, p2);
        String temp = incoming.substring(p2 + 1, p3);
        String lat  = incoming.substring(p3 + 1, p4);
        String lng  = incoming.substring(p4 + 1);

        // BUNGKUS MENJADI FORMAT JSON UNTUK WEBSITE
        String jsonPayload = "{";
        jsonPayload += "\"mid\":\"" + id + "\",";
        jsonPayload += "\"heart_rate\":" + toJsonNum(bpm) + ",";
        jsonPayload += "\"body_temperature\":" + toJsonNum(temp) + ",";
        jsonPayload += "\"latitude\":" + toJsonNum(lat) + ",";
        jsonPayload += "\"longitude\":" + toJsonNum(lng);
        jsonPayload += "}";

        // PUBLISH KE SERVER MQTT UBUNTU
        if(client.publish(mqtt_topic, jsonPayload.c_str())) {
           Serial.println("=> Berhasil kirim ke MQTT: " + jsonPayload);
        } else {
           Serial.println("=> GAGAL kirim ke MQTT!");
        }
      }
    }
  }
}