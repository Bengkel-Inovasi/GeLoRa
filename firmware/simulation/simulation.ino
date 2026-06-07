#include <WiFi.h>
#include <WiFiClientSecure.h>
#include <PubSubClient.h>
#include <WiFiManager.h>
#include <ArduinoJson.h>
#include "LittleFS.h"

// --- Variabel Global (Akan diisi dari Web Portal atau LittleFS) ---
char mqtt_server[60] = "03d11725d6454821a22ec23677b0463e.s1.eu.hivemq.cloud";
char mqtt_port[6]    = "8883";
char mqtt_user[40]   = "your_username";
char mqtt_pass[40]   = "your_password";
char mqtt_topic[50]  = "/server/record";

// --- Flag untuk simpan config ---
bool shouldSaveConfig = false;

WiFiClientSecure espClient;
PubSubClient client(espClient);

unsigned long lastMsg = 0;

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

  // Munculkan Portal "GeLoRa-Simulator-Config" jika tidak ada WiFi tersimpan
  if (!wm.autoConnect("GeLoRa-Simulator-Config")) {
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
    String clientId = "ESP32_Simulator_" + String(random(0xffff), HEX);
    
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
  
  loadConfigFile(); // Baca config dari memory flash (jika ada)
  setup_wifi();     // Konek WiFi & Portal Config
  
  client.setServer(mqtt_server, atoi(mqtt_port));
  Serial.println(F("=== SIMULATOR DATA PENDAKI SIAP (DYNAMIC CONFIG) ==="));
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

  unsigned long now = millis();
  // Simulasi pengiriman data setiap 10 detik
  if (now - lastMsg > 10000) {
    lastMsg = now;

    // SIMULASI DATA PENDAKI
    float heart_rate = random(65, 105);    // Detak jantung normal 65-105 bpm
    float temperature = random(355, 378) / 10.0; // Suhu normal 35.5 - 37.8 C
    float latitude = -7.7940 + (random(-200, 200) / 10000.0); // Sekitar koordinat simulasi
    float longitude = 110.3600 + (random(-200, 200) / 10000.0);

    // BUNGKUS MENJADI FORMAT JSON UNTUK BACKEND
    StaticJsonDocument<256> doc;
    doc["mid"] = "SIMULATOR_01";
    doc["heart_rate"] = heart_rate;
    doc["body_temperature"] = temperature;
    doc["latitude"] = latitude;
    doc["longitude"] = longitude;

    char jsonBuffer[256];
    serializeJson(doc, jsonBuffer);

    // PUBLISH KE SERVER MQTT
    if(client.publish(mqtt_topic, jsonBuffer)) {
       Serial.print("=> Terkirim ke MQTT: ");
       Serial.println(jsonBuffer);
    } else {
       Serial.println("=> GAGAL kirim ke MQTT!");
    }
  }
}
