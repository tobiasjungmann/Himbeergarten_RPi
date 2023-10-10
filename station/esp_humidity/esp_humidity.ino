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
WiFiClient client;

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
static sensor_t sensors[1] = { {0,859,481} }; 
#elif defined(ESP32)
static sensor_t sensors[1] = { {36,880,470} }; 
#endif
#define MAX_SENSOR_COUNT 16;

void setup() {
  Serial.begin(115200);
  delay(10);
  Serial.println('\n');
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  //Serial.println("Sleeping finished");
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


bool handle_sensors_callback(pb_istream_t *stream, const pb_field_t *field, void **arg) {
  smart_home_GetActiveSensorsReply *msg = (smart_home_GetActiveSensorsReply *)(*arg);
  uint32_t value;

  if (!pb_decode_varint32(stream, &value)) {
    return false;
  }

  // Append the value to the repeated field
  if (msg->sensorCount < 16) {
    //msg->sensors[msg->sensorCount++] = value;
    Serial.printf("Value parsed: %i\n",value);
   // msg->sensors.
    return true;
  }

  return false;
}

void getSensors() {
  smart_home_GetActiveSensorsRequest getSensorMsg = smart_home_GetActiveSensorsRequest_init_zero;
  strcpy(getSensorMsg.deviceMAC, WiFi.macAddress().c_str());
  uint8_t buffer[128];
  // todo potentially change to os_stream_from_socket https://github.com/nanopb/nanopb/blob/master/examples/network_server/client.c#L30
  pb_ostream_t stream = pb_ostream_from_buffer(buffer, 12);//sizeof(buffer));

  bool status = pb_encode(&stream, smart_home_GetActiveSensorsRequest_fields, &getSensorMsg);

  if (!status) {
    Serial.println("Failed to encode");
    return;
  }
  client.write(buffer, stream.bytes_written);
  String line = client.readStringUntil('\n');
  // recv proto message and parse back to its original value - nanopb anschauen

  Serial.println("Server response: " + line);
  smart_home_GetActiveSensorsReply reply = smart_home_GetActiveSensorsReply_init_default;
  pb_istream_t istream = pb_istream_from_buffer(reinterpret_cast<const unsigned char *>(line.c_str()), line.length());

if (!pb_decode_delimited(&istream, smart_home_GetActiveSensorsReply_fields, &reply))
        {
            fprintf(stderr, "Decode failed: %s\n", PB_GET_ERROR(&istream));
            //return false;
        }
 /* const pb_field_t reply_fields[] = {
    smart_home_GetActiveSensorsReply_sensors_tag,
     {
        pb_callback_t(smart_home_GetActiveSensorsReply, sensors, &handle_sensors_callback),
        smart_home_GetActiveSensorsReply_sensorCount_tag, {
            pb_encode(2, &getSensorMsg.sensorCount),
            1, false, 0, uint32_t,
        },
    },
    smart_home_GetActiveSensorsReply_sensorCount_tag,
    {
      pb_encode(2, &getSensorMsg.sensorCount),
      1,
      false,
      0,
      pb_type_uint32_t,
    },
  };

  if (!pb_decode(&stream, reply_fields, &getSensorMsg)) {
    return false;
  }*/

  // Serial.println(((uint8_t*)reply.sensors.arg)[0]);
  //reply.sensors
  /*const pb_field_t smart_home_GetActiveSensorsReply_fields[] = {
        // Other fields...
        smart_home_GetActiveSensorsReply_sensors_tag, {
            pb_callback_t {
                .funcs.decode = &handle_sensors_callback,
                .arg = &msg,
            },
            0, false, pb_membersize(smart_home_GetActiveSensorsReply, sensors), pb_type_uint8_t,
        },
        // More fields...
    };
     if (pb_decode(&istream, smart_home_GetActiveSensorsReply_fields, &reply)) {
    Serial.println("Parsing successful");//reply.sensors);
    };
  }*/
  client.stop();
}

void loop() {
  if (waitForConnection()) {
    //getSensors();   // reply must be decoded correctly with the list
    for(sensor_t s: sensors){
      getMoistureValues(s);
      strcpy(message.deviceMAC, WiFi.macAddress().c_str());
      message.sensorId = 0;   // todo why is this here?
      sendToForwarder();
    }
    /*for (int i = 0; i < sizeof(sensors); i++) {
      Serial.println("In loop");
      getMoistureValues(sensors[i]);
      strcpy(message.deviceMAC, WiFi.macAddress().c_str());
      message.sensorId = 0;
      sendToForwarder();
    }*/
  }
  if (DEEP_SLEEP) {
    Serial.println("Sleeping");
    delay(2000);
    ESP.deepSleep(SLEEP_TIME*60e6); 
  } else {
    Serial.println("Delay");
    delay(2000);
  }
}