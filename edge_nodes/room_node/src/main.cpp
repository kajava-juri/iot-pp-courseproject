/**
 * Room Node – main.cpp
 *
 * Reads temperature and humidity from a DHT22 sensor and publishes the
 * readings as JSON to the MQTT broker over Wi-Fi.
 *
 * MQTT topics published:
 *   room/temperature  – {"celsius": <°C>}
 *   room/humidity     – {"percent": <%rH>}
 */

#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <DHT.h>

// ── Configuration ─────────────────────────────────────────────────────────────
#define WIFI_SSID     "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"
#define MQTT_HOST     "192.168.1.100"   // IP of the host running docker-compose
#define MQTT_PORT     1883
#define MQTT_CLIENT   "room-node-01"
#define DHT_PIN       4
#define DHT_TYPE      DHT22
#define PUBLISH_INTERVAL_MS 5000

// ── Globals ───────────────────────────────────────────────────────────────────
WiFiClient   wifiClient;
PubSubClient mqtt(wifiClient);
DHT          dht(DHT_PIN, DHT_TYPE);

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
    while (!mqtt.connected()) {
        Serial.print("Connecting to MQTT broker...");
        if (mqtt.connect(MQTT_CLIENT)) {
            Serial.println(" connected.");
        } else {
            Serial.printf(" failed (rc=%d), retry in 5 s\n", mqtt.state());
            delay(5000);
        }
    }
}

// ── Setup ─────────────────────────────────────────────────────────────────────
void setup() {
    Serial.begin(115200);
    dht.begin();
    connectWifi();
    connectMqtt();
    Serial.println("Room node ready.");
}

// ── Loop ──────────────────────────────────────────────────────────────────────
void loop() {
    if (!mqtt.connected()) connectMqtt();
    mqtt.loop();

    float temperature = dht.readTemperature();
    float humidity    = dht.readHumidity();

    if (isnan(temperature) || isnan(humidity)) {
        Serial.println("DHT22 read failed – skipping publish.");
    } else {
        char buf[64];

        snprintf(buf, sizeof(buf), "{\"celsius\":%.1f}", temperature);
        mqtt.publish("room/temperature", buf);

        snprintf(buf, sizeof(buf), "{\"percent\":%.1f}", humidity);
        mqtt.publish("room/humidity", buf);

        Serial.printf("Published – temp: %.1f°C, humidity: %.1f%%\n",
                      temperature, humidity);
    }

    delay(PUBLISH_INTERVAL_MS);
}
