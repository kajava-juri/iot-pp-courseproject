#include <Arduino.h>
#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>
#include <ClickEncoder.h>
#include <ittiot.h>
#include <Ticker.h>
#include <Wire.h>

#include "iot_common.h"

// Defining pins as need for the encoder
#define ENC_PINA 12
#define ENC_PINB 13
#define ENC_BTN   0
#define ENC_STEPS_PER_NOTCH 4

Adafruit_MPU6050 mpu;
ClickEncoder encoder = ClickEncoder(ENC_PINA, ENC_PINB, ENC_BTN, ENC_STEPS_PER_NOTCH);

Ticker encTicker;

bool encFlag;
int16_t encoderValue;
int16_t lastEncoderValue;

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

void publishEncoderValue(int16_t value) {
  String msg = String(value);
  iot.publishMsg("enc", msg.c_str());
}

void setEncFlag() {
  encFlag = true;
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

  encoder.setButtonHeldEnabled(true);
  encoder.setDoubleClickEnabled(true);
  encTicker.attach(1, setEncFlag);

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

  static uint32_t lastService = 0;
  if (lastService + 1000 < micros()) {
    lastService = micros();
    encoder.service();
  }

  encoderValue += encoder.getValue();

  if (encoderValue != lastEncoderValue) {
    lastEncoderValue = encoderValue;
    Serial.print("Encoder Value: ");
    Serial.println(encoderValue);
  }

  if (encFlag) {
    encFlag = false;
    publishEncoderValue(encoderValue);
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