/**
 * Interface Node – main.cpp
 *
 * Acts as a Wi-Fi + MQTT gateway/bridge.  It subscribes to local (serial or
 * BLE) data streams from nearby sensor nodes and re-publishes the payloads to
 * the central MQTT broker over Wi-Fi.
 *
 * In this simplified implementation the node:
 *   • Maintains a persistent Wi-Fi and MQTT connection.
 *   • Echoes any data received on the serial port to the MQTT topic
 *     "interface/serial".
 *   • Subscribes to "interface/cmd" for over-the-air control commands.
 */

#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>

// ── Configuration ─────────────────────────────────────────────────────────────
#define WIFI_SSID     "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"
#define MQTT_HOST     "192.168.1.100"   // IP of the host running docker-compose
#define MQTT_PORT     1883
#define MQTT_CLIENT   "interface-node-01"

// ── Globals ───────────────────────────────────────────────────────────────────
WiFiClient   wifiClient;
PubSubClient mqtt(wifiClient);

// ── MQTT callback: handle inbound commands ────────────────────────────────────
void onMqttMessage(char* topic, byte* payload, unsigned int length) {
    String msg;
    for (unsigned int i = 0; i < length; i++) msg += (char)payload[i];
    Serial.printf("[MQTT] %s → %s\n", topic, msg.c_str());
}

// ── Helper: connect to Wi-Fi ──────────────────────────────────────────────────
void connectWifi() {
    Serial.printf("Connecting to Wi-Fi '%s'...", WIFI_SSID);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print('.');
    }
    Serial.printf("\nConnected – IP: %s\n", WiFi.localIP().toString().c_str());
}

// ── Helper: connect to MQTT broker ───────────────────────────────────────────
void connectMqtt() {
    mqtt.setServer(MQTT_HOST, MQTT_PORT);
    mqtt.setCallback(onMqttMessage);
    while (!mqtt.connected()) {
        Serial.print("Connecting to MQTT broker...");
        if (mqtt.connect(MQTT_CLIENT)) {
            Serial.println(" connected.");
            mqtt.subscribe("interface/cmd");
        } else {
            Serial.printf(" failed (rc=%d), retry in 5 s\n", mqtt.state());
            delay(5000);
        }
    }
}

// ── Setup ─────────────────────────────────────────────────────────────────────
void setup() {
    Serial.begin(115200);
    connectWifi();
    connectMqtt();
    Serial.println("Interface node ready.");
}

// ── Loop ──────────────────────────────────────────────────────────────────────
void loop() {
    if (!mqtt.connected()) connectMqtt();
    mqtt.loop();

    // Forward any bytes arriving on Serial to the broker
    if (Serial.available()) {
        String line = Serial.readStringUntil('\n');
        line.trim();
        if (line.length() > 0) {
            mqtt.publish("interface/serial", line.c_str());
            Serial.printf("[FWD] interface/serial ← %s\n", line.c_str());
        }
    }
}
