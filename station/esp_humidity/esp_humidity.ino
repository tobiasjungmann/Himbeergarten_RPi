//#include <ESP8266WiFi.h>
#include <WiFi.h>

#include "src/humidityStorage.pb.h"

#include "pb_common.h"
#include "pb.h"
#include "pb_encode.h"
#include "pb_decode.h"
//#include "nanopb.h"

/*  === Replace by actual values ===  */
const char *ssid = "aaaa";           // The SSID (name) of the Wi-Fi network you want to connect to
const char *password = "asdasdasd";  // The password of the Wi-Fi network
const char *addr = "192.168.0.4";    // Ip address of the forwarder interface
const uint16_t port = 12348;         // Port of the forwarder interface

// Source: https://how2electronics.com/interface-capacitive-soil-moisture-sensor-arduino/
/*  === Replace by specific values per sensor ===  */
const int AirValue = 880;    //the value which is returned when the sensor is dry
const int WaterValue = 470;  //the value which is returned when the sensor is submerged in water

const bool deepSleep = false;  // set to true if the device should deepsleep (Connect RST and D0)

smart_home_StoreHumidityRequest message = smart_home_StoreHumidityRequest_init_zero;
WiFiClient client;

static uint8_t sensors[] = { A0 };  // Sensor setupt D1 Mini
#define MAX_SENSOR_COUNT 16;

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
  // todo potentially cahnge to os_stream_from_socket https://github.com/nanopb/nanopb/blob/master/examples/network_server/client.c#L30
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