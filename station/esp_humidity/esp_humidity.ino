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

typedef struct {
  uint8_t id;
  int AirValue;    //the value which is returned when the sensor is dry
  int WaterValue;  //the value which is returned when the sensor is submerged in water
} sensor_t;





smart_home_StoreHumidityRequest message = smart_home_StoreHumidityRequest_init_zero;



void waitForConnection() {
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Connecting to WiFi...");
  }
  Serial.println("Connected to WiFi");
}

void setupWIFI() {
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  waitForConnection();
}

/*
 Source: https://how2electronics.com/interface-capacitive-soil-moisture-sensor-arduino/

const int AirValue    //the value which is returned when the sensor is dry
const int WaterValue  //the value which is returned when the sensor is submerged in water

Setup values for the specific sensors:
Sensor    Air   Water
1         880   470
2
3
4
5         859   481
*/
#if defined(ESP8266)
// todo change to struct with the predefined values
static sensor_t sensors[1] = { { 0, 859, 481 } };
#elif defined(ESP32)
static sensor_t sensors[1] = { { 36, 880, 470 } };
#endif
#define MAX_SENSOR_COUNT 16;

void setup() {
  Serial.begin(115200);

  setupWIFI();
}


void getMoistureValues(sensor_t sensor) {
  int soilMoistureValue = analogRead(sensor.id);
  Serial.print("Sensor: ");
  Serial.print(sensor.id);
  Serial.print("  Humidity Value: ");
  Serial.println(soilMoistureValue);
  int soilMoisturePercent = map(soilMoistureValue, sensor.AirValue, sensor.WaterValue, 0, 100);
  if (soilMoisturePercent > 100) {
    soilMoisturePercent = 100;
  } else if (soilMoisturePercent < 0) {
    soilMoisturePercent = 0;
  }
  message.humidityInPercent = soilMoisturePercent;
  message.humidity = soilMoistureValue;
}

void sendToForwarder() {
  uint8_t buffer[128];
  WiFiClient client;
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

  client.write(buffer, stream.bytes_written);
  Serial.println("");
}


/*void sendToForwarderREST() {
  uint8_t buffer[128];
  WiFiClient client;
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

  client.write(buffer, stream.bytes_written);
  Serial.println("");
}*/





void loop() {
 waitForConnection();
    //getSensors();   // reply must be decoded correctly with the list
    for (sensor_t s : sensors) {
      getMoistureValues(s);
      strcpy(message.deviceMAC, WiFi.macAddress().c_str());
      message.sensorId = 0;  // todo why is this here?
      sendToForwarder();
    }

  
  if (DEEP_SLEEP) {
    Serial.println("Sleeping");
    delay(2000);
    ESP.deepSleep(SLEEP_TIME * 60e6);
  } else {
    Serial.println("Delay");
    delay(2000);
  }
}