#include <Arduino.h>
#include <ittiot.h>

#define PIR_PIN D5
#define PIR_LED_PIN D4

#include "iot_common.h"

bool pirState = false;

void iot_connected() {
  Serial.println("MQTT connected");
}

void setup() {
  Serial.begin(115200);

  pinMode(PIR_PIN, INPUT);
  pinMode(PIR_LED_PIN, OUTPUT);

  iot.setConfig("wname", STR(WIFI_NAME));
  iot.setConfig("wpass", STR(WIFI_PASSWORD));
  iot.setConfig("msrv", "193.40.245.72");
  iot.setConfig("mport", "1883");
  iot.setConfig("muser", "test");
  iot.setConfig("mpass", "test");
  iot.printConfig(); // print IoT json config to serial
  iot.setup();       // Initialize IoT library
}

void loop() {
  iot.handle();

  if (digitalRead(PIR_PIN)) {
    if (!pirState) {
      digitalWrite(PIR_LED_PIN, HIGH);
      iot.publishMsg("event/motion", "1");
      pirState = true;
    }
  } else {
    if (pirState) {
      digitalWrite(PIR_LED_PIN, LOW);
      iot.publishMsg("event/motion", "0");
      pirState = false;
    }
  }

  delay(200);
}