#include <Arduino.h>

struct ProbePin {
  uint8_t pin;
  const char *name;
};

ProbePin probePins[] = {
    {D1, "D1 / GPIO5"},
    {D2, "D2 / GPIO4"},
    {D5, "D5 / GPIO14"},
    {D6, "D6 / GPIO12"},
    {D7, "D7 / GPIO13"},
};

constexpr size_t probePinCount = sizeof(probePins) / sizeof(probePins[0]);

void setup() {
  Serial.begin(115200);

  for (size_t i = 0; i < probePinCount; ++i) {
    pinMode(probePins[i].pin, OUTPUT);
    digitalWrite(probePins[i].pin, LOW);
  }

  Serial.println();
  Serial.println("Pin probe started. Measure each exposed pin against GND.");
}

void loop() {
  for (size_t activeIndex = 0; activeIndex < probePinCount; ++activeIndex) {
    for (size_t i = 0; i < probePinCount; ++i) {
      digitalWrite(probePins[i].pin, i == activeIndex ? HIGH : LOW);
    }

    Serial.print("Active pin: ");
    Serial.println(probePins[activeIndex].name);
    delay(2000);
  }
}