#include <ESP8266WiFi.h>

#include "humidityStorage.pb.h"

#include "pb_common.h"
#include "pb.h"
#include "pb_encode.h"
#include "pb_decode.h"




const char* ssid     = "aaaa";         // The SSID (name) of the Wi-Fi network you want to connect to
const char* password = "asdasdasd";     // The password of the Wi-Fi network
const char* addr     = "192.168.0.4";
const uint16_t port  = 12348;

WiFiClient client;



void setup() {
  Serial.begin(115200);         // Start the Serial communication to send messages to the computer
  delay(10);
  Serial.println('\n');
   WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);             // Connect to the network
  Serial.print("Connecting to ");
  Serial.print(ssid); Serial.println(" ...");

  int i = 0;
  while (WiFi.status() != WL_CONNECTED) { // Wait for the Wi-Fi to connect
    delay(1000);
    Serial.print(++i); Serial.print(' ');
  }

  Serial.println('\n');
  Serial.println("Connection established!");  
  Serial.print("IP address:\t");
  Serial.println(WiFi.localIP());         // Send the IP address of the ESP8266 to the computer


}

void loop() {
    if (!client.connect(addr, port)) {
    Serial.println("connection failed");
    Serial.println("wait 5 sec to reconnect...");
    delay(2000);
    return;
  }
    delay(2000);
    smart_home_StoreHumidityRequest message=smart_home_StoreHumidityRequest_init_zero;
    message.humidity=42;
    message.sensorId=0;
    message.sensorId=123;

  // GRPC stuff
  uint8_t buffer[128];
  pb_ostream_t stream = pb_ostream_from_buffer(buffer, sizeof(buffer));

   bool status = pb_encode(&stream, smart_home_StoreHumidityRequest_fields, &message);
 
if (!status)
{
    Serial.println("Failed to encode");
    return;
}
	
  Serial.printf("Amount of Bytes %d\n",stream.bytes_written);
  for(int i = 0; i<stream.bytes_written; i++){
    Serial.printf("%02X",buffer[i]);
  }

  client.write(buffer, stream.bytes_written);
  Serial.println("");
 }