#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include <WiFiManager.h>
#include <ArduinoJson.h>
#include "LittleFS.h"

// --- Variabel Global (Akan diisi dari Web Portal atau LittleFS) ---
char mqtt_server[60] = "your-cluster-id.s1.eu.hivemq.cloud";
char mqtt_port[6]    = "8883";
char mqtt_user[40]   = "your_username";
char mqtt_pass[40]   = "your_password";
char mqtt_topic[50]  = "/server/record";

// --- Flag untuk simpan config ---
bool shouldSaveConfig = false;

// --- KONFIGURASI PIN LORA UART (AS32/E32) ---
#define LORA_RX_PIN 26 
#define LORA_TX_PIN 27 
HardwareSerial loraSerial(2);

WiFiClientSecure espClient;
PubSubClient client(espClient);

// Callback saat user menekan "Save" di portal web WiFiManager
void saveConfigCallback() {
  Serial.println("Pengaturan baru diterima, akan disimpan...");
  shouldSaveConfig = true;
}

// --- Fungsi Simpan & Baca File Config di Memory Internal (LittleFS) ---
void saveConfigFile() {
  Serial.println("Saving config to LittleFS...");
  StaticJsonDocument<512> doc;
  doc["mqtt_server"] = mqtt_server;
  doc["mqtt_port"]   = mqtt_port;
  doc["mqtt_user"]   = mqtt_user;
  doc["mqtt_pass"]   = mqtt_pass;

  File configFile = LittleFS.open("/config.json", "w");
  if (!configFile) {
    Serial.println("Failed to open config file for writing");
    return;
  }
  serializeJson(doc, configFile);
  configFile.close();
  Serial.println("Config saved successfully.");
}

void loadConfigFile() {
  Serial.println("Mounting FS...");
  if (LittleFS.begin(true)) {
    Serial.println("Mounted file system");
    if (LittleFS.exists("/config.json")) {
      Serial.println("Reading config file...");
      File configFile = LittleFS.open("/config.json", "r");
      if (configFile) {
        size_t size = configFile.size();
        std::unique_ptr<char[]> buf(new char[size]);
        configFile.readBytes(buf.get(), size);
        StaticJsonDocument<512> doc;
        DeserializationError error = deserializeJson(doc, buf.get());
        if (!error) {
          Serial.println("Parsed JSON config");
          if (doc.containsKey("mqtt_server")) strcpy(mqtt_server, doc["mqtt_server"]);
          if (doc.containsKey("mqtt_port"))   strcpy(mqtt_port, doc["mqtt_port"]);
          if (doc.containsKey("mqtt_user"))   strcpy(mqtt_user, doc["mqtt_user"]);
          if (doc.containsKey("mqtt_pass"))   strcpy(mqtt_pass, doc["mqtt_pass"]);
        } else {
          Serial.println("Failed to load json config");
        }
        configFile.close();
      }
    }
  } else {
    Serial.println("Failed to mount FS");
  }
}

void setup_wifi() {
  WiFiManager wm;
  
  // Daftarkan callback untuk menyimpan config
  wm.setSaveConfigCallback(saveConfigCallback);

  // Buat input field kustom untuk MQTT di halaman WiFiManager
  WiFiManagerParameter custom_mqtt_server("server", "MQTT Server URL", mqtt_server, 60);
  WiFiManagerParameter custom_mqtt_port("port", "MQTT Port", mqtt_port, 6);
  WiFiManagerParameter custom_mqtt_user("user", "MQTT Username", mqtt_user, 40);
  WiFiManagerParameter custom_mqtt_pass("pass", "MQTT Password", mqtt_pass, 40);

  wm.addParameter(&custom_mqtt_server);
  wm.addParameter(&custom_mqtt_port);
  wm.addParameter(&custom_mqtt_user);
  wm.addParameter(&custom_mqtt_pass);

  // Munculkan Portal "GeLoRa-Basecamp-Config" jika tidak ada WiFi tersimpan
  if (!wm.autoConnect("GeLoRa-Basecamp-Config")) {
    Serial.println("Gagal konek ke WiFi, restart...");
    delay(3000);
    ESP.restart();
  }

  // Ambil nilai dari input field kustom setelah user mengonfigurasi via web
  strcpy(mqtt_server, custom_mqtt_server.getValue());
  strcpy(mqtt_port, custom_mqtt_port.getValue());
  strcpy(mqtt_user, custom_mqtt_user.getValue());
  strcpy(mqtt_pass, custom_mqtt_pass.getValue());

  // Simpan ke LittleFS jika data berubah
  if (shouldSaveConfig) {
    saveConfigFile();
  }

  Serial.println("\nWiFi Terhubung! IP: ");
  Serial.println(WiFi.localIP());
}

void reconnect() {
  while (!client.connected()) {
    Serial.printf("Mencoba konek ke MQTT %s:%s...", mqtt_server, mqtt_port);
    String clientId = "ESP32_Basecamp_" + String(random(0xffff), HEX);
    
    if (client.connect(clientId.c_str(), mqtt_user, mqtt_pass)) {
      Serial.println(" BERHASIL!");
    } else {
      Serial.print(" Gagal, rc=");
      Serial.print(client.state());
      Serial.println(" Coba lagi dalam 5 detik");
      delay(5000);
    }
  }
}

void setup() {
  Serial.begin(115200);
  
  // Izinkan koneksi TLS tanpa verifikasi sertifikat (untuk kemudahan prototype)
  espClient.setInsecure();
  
  // Inisialisasi LoRa UART AS32
  loraSerial.begin(9600, SERIAL_8N1, LORA_RX_PIN, LORA_TX_PIN);

  loadConfigFile(); // Baca config dari memory flash (jika ada)
  setup_wifi();     // Konek WiFi & Portal Config
  
  client.setServer(mqtt_server, atoi(mqtt_port));
  Serial.println(F("=== GATEWAY RX SIAP (AS32 UART + DYNAMIC CONFIG) ==="));
}

String toJsonNum(String val) {
  if (val.length() == 0) return "null";
  char c = val.charAt(0);
  if ((c >= '0' && c <= '9') || c == '-') return val;
  return "null";
}

void loop() {
  // Pastikan WiFi tetap konek
  if (WiFi.status() != WL_CONNECTED) {
    setup_wifi();
  }

  // Pastikan MQTT tetap konek
  if (!client.connected()) {
    reconnect();
  }
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

        // BUNGKUS MENJADI FORMAT JSON UNTUK BACKEND
        String jsonPayload = "{";
        jsonPayload += "\"mid\":\"" + id + "\",";
        jsonPayload += "\"heart_rate\":" + toJsonNum(bpm) + ",";
        jsonPayload += "\"body_temperature\":" + toJsonNum(temp) + ",";
        jsonPayload += "\"latitude\":" + toJsonNum(lat) + ",";
        jsonPayload += "\"longitude\":" + toJsonNum(lng);
        jsonPayload += "}";

        // PUBLISH KE SERVER MQTT (IP hasil konfigurasi web)
        if(client.publish(mqtt_topic, jsonPayload.c_str())) {
           Serial.println("=> Berhasil kirim ke MQTT: " + jsonPayload);
        } else {
           Serial.println("=> GAGAL kirim ke MQTT!");
        }
      }
    }
  }
}
