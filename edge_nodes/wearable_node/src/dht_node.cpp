// Includes global variables and librarys that the DHT uses
#include <Arduino.h>
#include <ittiot.h>
#include <Ticker.h>
#include <DHT.h>
 
#define WIFI_NAME "name"
#define WIFI_PASSWORD "password"
 
#define DHTPIN D3     // Pin where DHT shield is connected. Change this to D4 if the shield has no legs removed.
#define DHTTYPE DHT22   // DHT 22  (AM2302)
 
// Create an object for DHT sensor
DHT dht(DHTPIN, DHTTYPE);
 
// Create an object for Ticker library
Ticker timeTicker;
 
bool sendDataFlag;
 
// Ticker library callback, which will occur 0.5 second interval.
void sendData()
{
  sendDataFlag=true;
}
 
void iot_received(String topic, String msg) {}
 
// Function started after the connection to the server is established.
void iot_connected()
{
  // Send message to serial port to show that connection is established
  Serial.println("MQTT connected callback");
  // Send message to MQTT server to show that connection is established
  iot.log("IoT DHT example!");
}
 
void setup()
{
  // Initialize serial port and send message
  Serial.begin(115200); // setting up serial connection parameter
  Serial.println("Booting");
 
  //iot.setConfig("wname", WIFI_NAME);
  //iot.setConfig("wpass", WIFI_PASSWORD);
  iot.printConfig(); // print IoT json config to serial
  iot.setup(); // Initialize IoT library
 
  // Initialize DHT library
  dht.begin();
 
  // Initialize Ticker interval and callback
  timeTicker.attach(1, sendData);
}
 
void loop()
{
  iot.handle(); // IoT behind the plan work, it should be periodically called
 
  if(sendDataFlag)
  {
    sendDataFlag = false;
    // Read humidity and temperature
    float h = dht.readHumidity();
    float t = dht.readTemperature();
 
    // Create a buffer to store strings to being sent later
    char buf[10];
 
    // Convert temperature value messages to strings and send to the MQTT server
    String(t).toCharArray(buf,10);
    iot.publishMsg("temp",buf);
 
    // Convert humidity value messages to strings and send to the MQTT server
    String(h).toCharArray(buf,10);
    iot.publishMsg("hum",buf);
  }
}