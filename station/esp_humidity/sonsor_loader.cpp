#if defined(ESP8266)
#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#elif defined(ESP32)
#include <WiFi.h>
#include <HTTPClient.h>
#endif

//#include "httpRequest.h"
//#include "credentials.h"
#include "src/humidityStorage.pb.h"

#include "pb_common.h"
#include "pb.h"
#include "pb_encode.h"
#include "pb_decode.h"

/*
query the sensors which should be measured for this device
*/
void getSensors() {
    WiFiClient client;
  smart_home_GetActiveSensorsRequest getSensorMsg = smart_home_GetActiveSensorsRequest_init_zero;
  strcpy(getSensorMsg.deviceMAC, WiFi.macAddress().c_str());
  uint8_t buffer[128];
  // todo potentially change to os_stream_from_socket https://github.com/nanopb/nanopb/blob/master/examples/network_server/client.c#L30
  pb_ostream_t stream = pb_ostream_from_buffer(buffer, 12);  //sizeof(buffer));

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

  if (!pb_decode_delimited(&istream, smart_home_GetActiveSensorsReply_fields, &reply)) {
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

bool handle_sensors_callback(pb_istream_t *stream, const pb_field_t *field, void **arg) {
  smart_home_GetActiveSensorsReply *msg = (smart_home_GetActiveSensorsReply *)(*arg);
  uint32_t value;

  if (!pb_decode_varint32(stream, &value)) {
    return false;
  }

  // Append the value to the repeated field
  if (msg->sensorCount < 16) {
    //msg->sensors[msg->sensorCount++] = value;
    Serial.printf("Value parsed: %i\n", value);
    // msg->sensors.
    return true;
  }

  return false;
}