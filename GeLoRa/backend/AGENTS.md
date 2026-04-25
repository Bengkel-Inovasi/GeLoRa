this backend was made to support a research project called GeLoRa.
GeLoRa or "Gelang LoRa" is a device that act as tracker that will be used for mountain climbers to track their location and medical situation such as heart rate and body temp. GeLoRa can be used as wrist band
GeLoRa consist of two individuals device, first device act as wrist band that will be used by mountain climbers itself, this device consist of ESP32-C3, MAX30102, MLX90614, BN220, and LoRa to support wireless data transmission.
Second device is a device that will be permanently installed on mountain basecamp.
This backend will be used to act as communicatiom bridge from wristband to receiver and display the data to website.
MQTT is strongly recommended.
GeLoRa wristband device (mountain climbers side) can be extended. 