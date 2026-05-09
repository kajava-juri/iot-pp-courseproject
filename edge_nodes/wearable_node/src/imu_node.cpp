#include <Arduino.h>
#include <Adafruit_MPU6050.h>
#include <Adafruit_Sensor.h>
#include <ittiot.h>
#include <Ticker.h>
#include <Wire.h>
#include <ClickEncoder.h>

#include "iot_common.h"

// Control if sensor is enabled (connected)
#define ENCODER_ENABLE 1
#define IMU_ENABLE 1

// Defining pins for encoder
#define ENC_PINA 12
#define ENC_PINB 13
#define ENC_BTN 0
#define ENC_STEPS_PER_NOTCH 4
#define ENCODER_THRESHOLD 5  // threshold for rate of change detection
#define PRINT_INTERVAL_MS 40
// Fall detection windows
#define PEAK_WINDOW 5    // short window for local peak (samples)
#define RMS_WINDOW 100   // longer window for RMS baseline (samples)
#define REFRACTORY_MS 2000

Adafruit_MPU6050 mpu;
ClickEncoder encoder(ENC_PINA, ENC_PINB, ENC_BTN, ENC_STEPS_PER_NOTCH);

Ticker encTicker;

bool encFlag = false;
int16_t lastEncoderValue = 0;
int16_t lastEncoderDelta = 0;
unsigned long lastAccelPrintTime = 0;
// Sliding-max deque (small window)
float peakVal[PEAK_WINDOW];
unsigned long peakIdx[PEAK_WINDOW];
int peakCount = 0;

// Sliding RMS buffer (long window)
float rmsBuf[RMS_WINDOW];
int rmsIndex = 0;
int rmsCount = 0;
float sumSq = 0.0;

unsigned long sampleCounter = 0;
unsigned long lastTriggerMs = 0;

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

  // Initialize Encoder
  #if ENCODER_ENABLE
  encoder.setButtonHeldEnabled(true);
  encoder.setDoubleClickEnabled(true);
  encTicker.attach(1, setEncFlag);
  #endif

  #if IMU_ENABLE
  if (!mpu.begin()) {
    Serial.println("Failed to find MPU6050 chip");
    while (1) {
      delay(10);
    }
  }
  Serial.println("MPU6050 Found!");

  mpu.setHighPassFilter(MPU6050_HIGHPASS_0_63_HZ);
  mpu.setMotionDetectionThreshold(1);
  mpu.setMotionDetectionDuration(2);
  mpu.setInterruptPinLatch(true);
  mpu.setInterruptPinPolarity(true);
  mpu.setMotionInterrupt(true);
  #endif

  Serial.println("");
  delay(100);
}

void loop() {
  iot.handle();

  delay(10);

  #if ENCODER_ENABLE
  // Service the encoder
  static uint32_t lastService = 0;
  if (lastService + 1000 < micros()) {
    lastService = micros();
    encoder.service();
  }
  
  // Read encoder value
  static int16_t encoderValue = 0;
  encoderValue += encoder.getValue();
  
  // Detect rapid changes in encoder value
  if (encFlag) {
    encFlag = false;
    int16_t currentDelta = encoderValue - lastEncoderValue;
    int16_t deltaChange = currentDelta - lastEncoderDelta;  // acceleration (change in rate)
    
    Serial.print("Encoder Value: ");
    Serial.print(encoderValue);
    Serial.print(", Delta: ");
    Serial.print(currentDelta);
    Serial.print(", Acceleration: ");
    Serial.println(deltaChange);
    
    // Publish encoder value
    char buf[16];
    String(encoderValue).toCharArray(buf, sizeof(buf));
    iot.publishMsg("enc", buf);
    
    // Check if rate of change exceeds threshold
    if (abs(currentDelta) > ENCODER_THRESHOLD) {
      Serial.println("Rapid encoder change detected! Value: " + String(currentDelta));
      String msg = String(currentDelta);
      iot.publishMsg("event/encoder_rapid", msg.c_str());
    }
    
    lastEncoderValue = encoderValue;
    lastEncoderDelta = currentDelta;
  }
  #endif

  #if IMU_ENABLE
  if (mpu.getMotionInterruptStatus()) {
    sensors_event_t a, g, temp;
    mpu.getEvent(&a, &g, &temp);

    unsigned long currentTime = millis();
    if (currentTime - lastAccelPrintTime >= PRINT_INTERVAL_MS) {
      Serial.print("Accel:");
      Serial.print(sqrt(a.acceleration.x * a.acceleration.x + a.acceleration.y * a.acceleration.y + a.acceleration.z * a.acceleration.z));
      Serial.print(", ");
      Serial.print("Peak: ");
      Serial.print(peakVal[0]);
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
    }

    float s = sqrt(a.acceleration.x * a.acceleration.x + a.acceleration.y * a.acceleration.y + a.acceleration.z * a.acceleration.z);
    sampleCounter++;


    // subtract old value if buffer is full
    if (rmsCount == RMS_WINDOW) {
      float old = rmsBuf[rmsIndex];
      sumSq -= old * old;
    } else {
      rmsCount++;
    }
    rmsBuf[rmsIndex] = s;
    sumSq += s * s;
    rmsIndex = (rmsIndex + 1) % RMS_WINDOW;
    float baselineRMS = sqrt(sumSq / (rmsCount > 0 ? rmsCount : 1));

    // Now the peak detection using up to date local maximum
    // for every sample, remove smaller values from the back (small values can no longer be peaks)
    while (peakCount > 0 && peakVal[peakCount - 1] <= s) {
      peakCount--;
    }

    // append new value
    if (peakCount < PEAK_WINDOW) {
      peakVal[peakCount] = s;
      peakIdx[peakCount] = sampleCounter;
      peakCount++;
    }
    // remove outdated entries (older than PEAK_WINDOW samples)
    while (peakCount > 0 && sampleCounter - peakIdx[0] >= (unsigned long)PEAK_WINDOW) {
      // shift left
      for (int i = 0; i < peakCount - 1; ++i) {
        // start from left and replace the i-th element with the next one
        peakVal[i] = peakVal[i + 1];
        peakIdx[i] = peakIdx[i + 1];
      }
      peakCount--;
    }

    // Now take the current peak and compare with baseline to detect falls
    float peak = (peakCount > 0) ? peakVal[0] : s;
    if (currentTime - lastTriggerMs > REFRACTORY_MS) {
      const float R_RATIO = 2.0; // peak must be > R_RATIO * baseline
      const float DELTA_G = 1.5; // absolute difference in g
      if (peak > baselineRMS * R_RATIO && (peak - baselineRMS) > DELTA_G) {
        Serial.println("Fall detected! Peak: " + String(peak) + ", Baseline RMS: " + String(baselineRMS));
        char buf[32];
        String(peak, 2).toCharArray(buf, sizeof(buf));
        iot.publishMsg("event/fall", buf);
        publishReading("event/fall/rms", baselineRMS);
        lastTriggerMs = currentTime;
      }
    }

    // publishReading("imu/accel/x", a.acceleration.x);
    // publishReading("imu/accel/y", a.acceleration.y);
    // publishReading("imu/accel/z", a.acceleration.z);
    // publishReading("imu/gyro/x", g.gyro.x);
    // publishReading("imu/gyro/y", g.gyro.y);
    // publishReading("imu/gyro/z", g.gyro.z);
  }
  #endif

  delay(10);
}