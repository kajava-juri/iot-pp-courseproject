#include <Arduino.h>
#include <ittiot.h>

#define PIR_PIN D5
#define PIR_LED_PIN D4

bool pirState = false;

void iot_connected() {
  Serial.println("MQTT connected");
}

void setup() {
  Serial.begin(115200);

  iot.setup();

  pinMode(PIR_PIN, INPUT);
  pinMode(PIR_LED_PIN, OUTPUT);
}

void loop() {
  iot.handle();

  if (digitalRead(PIR_PIN)) {
    if (!pirState) {
      digitalWrite(PIR_LED_PIN, HIGH);
      iot.publishMsg("sensor/pir", "1");
      pirState = true;
    }
  } else {
    if (pirState) {
      digitalWrite(PIR_LED_PIN, LOW);
      iot.publishMsg("sensor/pir", "0");
      pirState = false;
    }
  }

  delay(200);
}