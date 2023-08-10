//#include <ESP8266WiFi.h>
#include <WiFi.h>

#include "src/humidityStorage.pb.h"

#include "pb_common.h"
#include "pb.h"
#include "pb_encode.h"
#include "pb_decode.h"

/*  === Replace by actual values ===  */
const char* ssid = "aaaa";           // The SSID (name) of the Wi-Fi network you want to connect to
const char* password = "asdasdasd";  // The password of the Wi-Fi network
const char* addr = "192.168.0.4";    // Ip address of the forwarder interface
const uint16_t port = 12348;         // Port of the forwarder interface

// Source: https://how2electronics.com/interface-capacitive-soil-moisture-sensor-arduino/
/*  === Replace by specific values per sensor ===  */
const int AirValue = 880;    //the value which is returned when the sensor is dry
const int WaterValue = 470;  //the value which is returned when the sensor is submerged in water

const bool deepSleep = false;  // set to true if the device should deepsleep (Connect RST and D0)

smart_home_StoreHumidityRequest message = smart_home_StoreHumidityRequest_init_zero;
WiFiClient client;

static uint8_t sensors[] = {A0};  // Sensor setupt D1 Mini

void setup() {
  Serial.begin(115200);
  delay(10);
  Serial.println('\n');
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  Serial.println("going to sleep");

  Serial.println("Sleeping finished");
}


void getMoistureValues(uint8_t sensor) {
  int soilMoistureValue = analogRead(sensor);
  Serial.println(soilMoistureValue);
  int soilMoisturePercent = map(soilMoistureValue, AirValue, WaterValue, 0, 100);
  if (soilMoisturePercent > 100) {
    soilMoisturePercent = 100;
  } else if (soilMoisturePercent < 0) {
    soilMoisturePercent = 0;
  }
  message.humidityInPercent = soilMoisturePercent;
  message.humidity = soilMoistureValue;
}

bool waitForConnection() {
  Serial.print("Connecting to ");
  Serial.print(ssid);
  int i = 0;
  while (WiFi.status() != WL_CONNECTED) {  // Wait for the Wi-Fi to connect
    delay(1000);
    Serial.print(++i);
    Serial.print(' ');
  }

  Serial.println('\n');
  Serial.println("Connection established!");
  Serial.print("IP address:\t");
  Serial.println(WiFi.localIP());
  Serial.println(WiFi.macAddress());

  for (i = 0; i < 15; i++) {
    if (!client.connect(addr, port)) {
      Serial.println("connection to server failed waiting to reconnect...");
    } else {
      return true;
    }
    delay(1000);
  }
  return false;
}

void sendToForwarder() {
  uint8_t buffer[128];
  pb_ostream_t stream = pb_ostream_from_buffer(buffer, sizeof(buffer));

  bool status = pb_encode(&stream, smart_home_StoreHumidityRequest_fields, &message);

  if (!status) {
    Serial.println("Failed to encode");
    return;
  }

  Serial.printf("Amount of Bytes %d\n", stream.bytes_written);
  for (int i = 0; i < stream.bytes_written; i++) {
    Serial.printf("%02X", buffer[i]);
  }

  client.write(buffer, stream.bytes_written);
  Serial.println("");
}

void getSensors(){
  smart_home_GetActiveSensorsRequest getSensorMsg = smart_home_GetActiveSensorsRequest_init_zero;
   strcpy(getSensorMsg.deviceMAC, WiFi.macAddress().c_str());
   uint8_t buffer[128];
  pb_ostream_t stream = pb_ostream_from_buffer(buffer, sizeof(buffer));

  bool status = pb_encode(&stream, smart_home_GetActiveSensorsRequest_fields, &getSensorMsg);

  if (!status) {
    Serial.println("Failed to encode");
    return;
  }

  /*Serial.printf("Amount of Bytes %d\n", stream.bytes_written);
  for (int i = 0; i < stream.bytes_written; i++) {
    Serial.printf("%02X", buffer[i]);
  }*/

  client.write(buffer, stream.bytes_written);
   String serverResponse = "";
  while (client.available()) {
    Serial.println("Reading one line from intput...");
    String line = client.readStringUntil('\n');
    serverResponse += line;
  }
  Serial.println("Server response: " + serverResponse);
 client.stop();
}

void loop() {
  if (waitForConnection()) {
    getSensors();
    /*for (int i = 0; i < sizeof(sensors); i++) {
      getMoistureValues(sensors[i]);
      strcpy(message.deviceMAC, WiFi.macAddress().c_str());
      message.sensorId = 0;
      sendToForwarder();
    }*/
  }
  if (deepSleep) {
    ESP.deepSleep(60e9);  // 1h
    yield();
  } else {
    delay(2000);
  }
}