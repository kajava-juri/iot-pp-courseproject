/**
 * Wearable Node – main.cpp
 *
 * Reads data from an MPU6050 IMU (accelerometer + gyroscope) and a simple
 * heart-rate/pulse sensor, then publishes the readings as JSON to the MQTT
 * broker over Wi-Fi.
 *
 * MQTT topics published:
 *   wearable/accel   – {"x": <g>, "y": <g>, "z": <g>}
 *   wearable/gyro    – {"x": <°/s>, "y": <°/s>, "z": <°/s>}
 *   wearable/heart   – {"bpm": <beats/min>}
 */

#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>
#include <Wire.h>

// ── Configuration ─────────────────────────────────────────────────────────────
#define WIFI_SSID     "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"
#define MQTT_HOST     "192.168.1.100"   // IP of the host running docker-compose
#define MQTT_PORT     1883
#define MQTT_CLIENT   "wearable-node-01"
#define HEART_PIN     34                // Analog pin for pulse sensor
#define PUBLISH_INTERVAL_MS 2000

// ── Globals ───────────────────────────────────────────────────────────────────
WiFiClient    wifiClient;
PubSubClient  mqtt(wifiClient);
Adafruit_MPU6050 mpu;

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

    connectWifi();
    connectMqtt();

    if (!mpu.begin()) {
        Serial.println("MPU6050 not found – check wiring!");
        while (true) delay(100);
    }
    mpu.setAccelerometerRange(MPU6050_RANGE_8_G);
    mpu.setGyroRange(MPU6050_RANGE_500_DEG);
    mpu.setFilterBandwidth(MPU6050_BAND_21_HZ);

    Serial.println("Wearable node ready.");
}

// ── Loop ──────────────────────────────────────────────────────────────────────
void loop() {
    if (!mqtt.connected()) connectMqtt();
    mqtt.loop();

    sensors_event_t accel, gyro, temp;
    mpu.getEvent(&accel, &gyro, &temp);

    // Publish accelerometer
    char buf[128];
    snprintf(buf, sizeof(buf),
             "{\"x\":%.3f,\"y\":%.3f,\"z\":%.3f}",
             accel.acceleration.x,
             accel.acceleration.y,
             accel.acceleration.z);
    mqtt.publish("wearable/accel", buf);

    // Publish gyroscope
    snprintf(buf, sizeof(buf),
             "{\"x\":%.3f,\"y\":%.3f,\"z\":%.3f}",
             gyro.gyro.x, gyro.gyro.y, gyro.gyro.z);
    mqtt.publish("wearable/gyro", buf);

    // Publish (mock) heart rate from analog sensor
    int raw = analogRead(HEART_PIN);
    int bpm = map(raw, 0, 4095, 40, 200);
    snprintf(buf, sizeof(buf), "{\"bpm\":%d}", bpm);
    mqtt.publish("wearable/heart", buf);

    delay(PUBLISH_INTERVAL_MS);
}
