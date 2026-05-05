#include <Arduino.h>
#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>
#include <ittiot.h>
#include <Ticker.h>
#include <Wire.h>
#include <DHT.h>

#include "iot_common.h"

// Defining pins for DHT sensor
#define DHTPIN D3
#define DHTTYPE DHT22

Adafruit_MPU6050 mpu;
DHT dht(DHTPIN, DHTTYPE);

Ticker dhtTicker;

bool sendDHTFlag;

void iot_received(String topic, String msg) {}

void iot_connected() {
  Serial.println("MQTT connected callback");
  iot.log("IoT IMU example!");
}

void publishReading(const char *topic, float value) {
  char buf[16];
  String(value, 4).toCharArray(buf, sizeof(buf));
  iot.publishMsg(topic, buf);
}

void sendDHT() {
  sendDHTFlag = true;
}

void setup(void) {
  Serial.begin(115200);
  while (!Serial)
    delay(10);

  Serial.println("Adafruit MPU6050 test!");

  iot.setConfig("wname", STR(WIFI_NAME));
  iot.setConfig("wpass", STR(WIFI_PASSWORD));
  iot.setConfig("msrv", "193.40.245.72");
  iot.setConfig("mport", "1883");
  iot.setConfig("muser", "test");
  iot.setConfig("mpass", "test");
  iot.printConfig();
  iot.setup();

  // Initialize DHT sensor
  dht.begin();
  dhtTicker.attach(1, sendDHT);

  if (!mpu.begin()) {
    Serial.println("Failed to find MPU6050 chip");
    while (1) {
      delay(10);
    }
  }
  Serial.println("MPU6050 Found!");

  mpu.setHighPassFilter(MPU6050_HIGHPASS_0_63_HZ);
  mpu.setMotionDetectionThreshold(1);
  mpu.setMotionDetectionDuration(20);
  mpu.setInterruptPinLatch(true);
  mpu.setInterruptPinPolarity(true);
  mpu.setMotionInterrupt(true);

  Serial.println("");
  delay(100);
}

void loop() {
  iot.handle();

  delay(10);

  if (sendDHTFlag) {
    sendDHTFlag = false;
    // Read humidity and temperature
    float h = dht.readHumidity();
    float t = dht.readTemperature();

    publishReading("temp", t);
    publishReading("hum", h);
  }

  if (mpu.getMotionInterruptStatus()) {
    sensors_event_t a, g, temp;
    mpu.getEvent(&a, &g, &temp);

    Serial.print("AccelX:");
    Serial.print(a.acceleration.x);
    Serial.print(",");
    Serial.print("AccelY:");
    Serial.print(a.acceleration.y);
    Serial.print(",");
    Serial.print("AccelZ:");
    Serial.print(a.acceleration.z);
    Serial.print(", ");
    Serial.print("GyroX:");
    Serial.print(g.gyro.x);
    Serial.print(",");
    Serial.print("GyroY:");
    Serial.print(g.gyro.y);
    Serial.print(",");
    Serial.print("GyroZ:");
    Serial.print(g.gyro.z);
    Serial.println("");

    publishReading("imu/accel/x", a.acceleration.x);
    publishReading("imu/accel/y", a.acceleration.y);
    publishReading("imu/accel/z", a.acceleration.z);
    publishReading("imu/gyro/x", g.gyro.x);
    publishReading("imu/gyro/y", g.gyro.y);
    publishReading("imu/gyro/z", g.gyro.z);
  }

  delay(10);
}