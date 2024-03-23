#if defined(ESP8266)
#include <ESP8266WiFi.h>
#elif defined(ESP32)
#include <WiFi.h>
#endif

#include "src/humidityStorage.pb.h"
#include "credentials.h"

#include "pb_common.h"
#include "pb.h"
#include "pb_encode.h"
#include "pb_decode.h"

smart_home_StoreHumidityRequest message = smart_home_StoreHumidityRequest_init_zero;

void setup() {
  Serial.begin(115200);

 WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
}

/**
*Read the current state of the sensor and write it into the struct
*/
void readSensorValues(sensor_t sensor) {
  int soilMoistureValue = analogRead(sensor.gpio);
  Serial.print("Sensor: ");
  Serial.print(sensor.gpio);
  Serial.print("  Humidity Value: ");
  Serial.println(soilMoistureValue);
  int soilMoisturePercent = map(soilMoistureValue, sensor.AirValue, sensor.WaterValue, 0, 100);
  if (soilMoisturePercent > 100) {
    soilMoisturePercent = 100;
  } else if (soilMoisturePercent < 0) {
    soilMoisturePercent = 0;
  }
  message.sensorId = sensor.label;
  message.humidityInPercent = soilMoisturePercent;
  message.humidity = soilMoistureValue;
}

void waitForConnection() {
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Connecting to WiFi...");
  }
  Serial.println("Connected to WiFi");
}

bool establishConnectionToForwarder(WiFiClient& client){
  waitForConnection();
  for (int i = 0; i < 15; i++) {
    if (!client.connect(FORWARDER_IP, FORWARDER_PORT)) {
      Serial.println("Connecting to server...");
      delay(1000);
    } else {
      Serial.println("Connected to server");
      return true;
    }
  }
  return false;
}

void sendToForwarder() {
  uint8_t buffer[128];
  WiFiClient client;
  establishConnectionToForwarder(client);

  pb_ostream_t stream = pb_ostream_from_buffer(buffer, sizeof(buffer));

  bool status = pb_encode(&stream, smart_home_StoreHumidityRequest_fields, &message);

  if (!status) {
    Serial.println("Failed to encode.");
    return;
  }

  Serial.printf("Amount of Bytes %d\n", stream.bytes_written);
  for (int i = 0; i < stream.bytes_written; i++) {
    Serial.printf("%02X", buffer[i]);
  }
  Serial.println("");

  client.write(buffer, stream.bytes_written);
}

void loop() {
  waitForConnection();
  //getSensors();   // reply must be decoded correctly with the list
  for (sensor_t s : sensors) {
    readSensorValues(s);
    strcpy(message.deviceMAC, WiFi.macAddress().c_str());
    sendToForwarder();
  }

  if (DEEP_SLEEP) {
    Serial.println("Sleeping");
    delay(2000);
    ESP.deepSleep(SLEEP_TIME * 60e6);
  } else {
    Serial.println("Delay");
    delay(10000);
  }
}