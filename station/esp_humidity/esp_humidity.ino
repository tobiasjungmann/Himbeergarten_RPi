#include <ESP8266WiFi.h>

#include "humidityStorage.pb.h"

#include "pb_common.h"
#include "pb.h"
#include "pb_encode.h"
#include "pb_decode.h"


const char* ssid = "aaaa";           // The SSID (name) of the Wi-Fi network you want to connect to
const char* password = "asdasdasd";  // The password of the Wi-Fi network
const char* addr = "192.168.0.4";
const uint16_t port = 12348;

WiFiClient client;


void setup() {
  Serial.begin(115200);
  delay(10);
  Serial.println('\n');
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
}

// Source: https://how2electronics.com/interface-capacitive-soil-moisture-sensor-arduino/
const int AirValue = 880;    //the value which is returned when the sensor is dry
const int WaterValue = 470;  //the value which is returned when the sensor is submerged in water
smart_home_StoreHumidityRequest message = smart_home_StoreHumidityRequest_init_zero;


void getMoistureValues() {
  int soilMoistureValue = analogRead(A0);
  Serial.println(soilMoistureValue);
  int soilmoisturepercent = map(soilMoistureValue, AirValue, WaterValue, 0, 100);
  if (soilmoisturepercent > 100) {
    soilmoisturepercent=100;
  } else if (soilmoisturepercent < 0) {
    soilmoisturepercent=0;
  } 
  message.humidityInPercent=soilmoisturepercent;
  message.humidity=soilMoistureValue;
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

  if (!client.connect(addr, port)) {
    Serial.println("connection failed");
    Serial.println("wait 5 sec to reconnect...");

    return false;
  }
  return true
}

void sendToForwarder(){
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


void loop() {
  if (waitForConnection()) {
    getMoistureValues();
    //strcpy(message.deviceId, WiFi.macAddress());
    message.sensorId = 0;
    sendToForwarder();
  }
  delay(2000);
}