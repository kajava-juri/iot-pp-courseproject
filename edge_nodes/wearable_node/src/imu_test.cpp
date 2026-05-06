#include <Arduino.h>
#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>
#include <Ticker.h>
#include <Wire.h>

#define PRINT_INTERVAL_MS 5

Adafruit_MPU6050 mpu;

unsigned long lastAccelPrintTime = 0;

void setup(void) {
  Serial.begin(115200);
  while (!Serial)
    delay(10);

  Serial.println("Adafruit MPU6050 test!");

//   iot.setConfig("wname", STR(WIFI_NAME));
//   iot.setConfig("wpass", STR(WIFI_PASSWORD));
//   iot.setConfig("msrv", "193.40.245.72");
//   iot.setConfig("mport", "1883");
//   iot.setConfig("muser", "test");
//   iot.setConfig("mpass", "test");
//   iot.printConfig();
//   iot.setup();

  if (!mpu.begin()) {
    Serial.println("Failed to find MPU6050 chip");
    while (1) {
      delay(10);
    }
  }
  Serial.println("MPU6050 Found!");

  mpu.setHighPassFilter(MPU6050_HIGHPASS_0_63_HZ);
  mpu.setMotionDetectionThreshold(1);
  mpu.setMotionDetectionDuration(1);
  mpu.setInterruptPinLatch(true);
  mpu.setInterruptPinPolarity(true);
  mpu.setMotionInterrupt(true);

  Serial.println("");
  delay(100);
}

void loop() {

  if (mpu.getMotionInterruptStatus()) {
    sensors_event_t a, g, temp;
    mpu.getEvent(&a, &g, &temp);

    unsigned long currentTime = millis();
    if (currentTime - lastAccelPrintTime >= PRINT_INTERVAL_MS) {
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

      lastAccelPrintTime = currentTime;
    }
  }

  delay(10);
}