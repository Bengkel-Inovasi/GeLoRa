#!/bin/bash
set -e

echo "=============================="
echo "   GeLoRa System Startup"
echo "=============================="

# --- 1. MQTT Broker ---
echo ""
echo "[1/2] Starting MQTT Broker..."
sudo pkill mosquitto 2>/dev/null || true
sleep 1
sudo nohup mosquitto -c /etc/mosquitto/mosquitto.conf > /tmp/mosquitto.log 2>&1 &
sleep 2

if pgrep -x mosquitto > /dev/null; then
    echo "      OK - Mosquitto running"
else
    echo "      FAILED - Check /tmp/mosquitto.log"
    cat /tmp/mosquitto.log
    exit 1
fi

# --- 2. Show WSL IP ---
WSL_IP=$(ip addr show eth2 2>/dev/null | grep "inet " | awk '{print $2}' | cut -d/ -f1)
if [ -z "$WSL_IP" ]; then
    WSL_IP=$(ip addr show | grep "inet " | grep -v "127.0.0.1" | grep -v "10.255.255" | grep -v "169.254" | awk '{print $2}' | cut -d/ -f1 | tail -1)
fi

echo ""
echo "[2/2] Network:"
echo "      WSL IP  : $WSL_IP"
echo "      MQTT    : $WSL_IP:1883"

# Check if IP matches what's in Arduino file
ARDUINO_IP=$(grep 'mqtt_server' /home/rozan/GeLoRa/inoexample/mqtt-example-rx.cpp | grep -oP '"\K[^"]+(?="[^"]*//)')
if [ "$ARDUINO_IP" != "$WSL_IP" ]; then
    echo ""
    echo "  *** WARNING: Arduino IP mismatch! ***"
    echo "      Arduino has : $ARDUINO_IP"
    echo "      Current IP  : $WSL_IP"
    echo "      Re-flash the ESP32 RX with the new IP."
fi

echo ""
echo "=============================="
echo "  Open 2 more terminals:"
echo ""
echo "  BACKEND:"
echo "  cd ~/GeLoRa/GeLoRa/backend"
echo "  go run ./cmd/core/main.go"
echo ""
echo "  FRONTEND:"
echo "  cd ~/GeLoRa/frontend"
echo "  npm run dev"
echo ""
echo "  Browser: http://localhost:3000"
echo "=============================="
