#include <dht.h>

dht DHT;

#define DHT11_PIN 2


const int numberOfAnalogPorts = 8;
void setup() {
  Serial.begin(9600); // open serial port, set the baud rate as 9600 bps
    Serial.println("REBOOT");
  if (pressure.begin())
    Serial.println("BMP180 init success");
  else
  {
    Serial.println("BMP180 init fail\n\n");
    while(1); // Pause forever.
  }
}

void print_all_analog() {
  Serial.print("a ");
  for(int i=0;i<numberOfAnalogPorts;i++){
    Serial.print(analogRead(i));
    Serial.print(" ");
  }
}

void print_temp_preassure()
{
  Serial.print("tp ");
  char status;
  double T,P,p0,a;

  status = pressure.startTemperature();
  if (status != 0)
  {
    // Wait for the measurement to complete:
    delay(status);

    // Retrieve the completed temperature measurement:
    // Note that the measurement is stored in the variable T.
    // Function returns 1 if successful, 0 if failure.

    status = pressure.getTemperature(T);
    if (status != 0)
    {
      // Print out the measurement:
      Serial.print("temperature: ");
      Serial.print(T,2);
      Serial.print(" deg C, ");
      Serial.print((9.0/5.0)*T+32.0,2);
      Serial.println(" deg F");

      // Start a pressure measurement:
      // The parameter is the oversampling setting, from 0 to 3 (highest res, longest wait).
      // If request is successful, the number of ms to wait is returned.
      // If request is unsuccessful, 0 is returned.

      status = pressure.startPressure(3);
      if (status != 0)
      {
        // Wait for the measurement to complete:
        delay(status);

        // Retrieve the completed pressure measurement:
        // Note that the measurement is stored in the variable P.
        // Note also that the function requires the previous temperature measurement (T).
        // (If temperature is stable, you can do one temperature measurement for a number of pressure measurements.)
        // Function returns 1 if successful, 0 if failure.

        status = pressure.getPressure(P,T);
        if (status != 0)
        {
          // Print out the measurement:
          Serial.print("absolute pressure: ");
          Serial.print(P,2);
          Serial.print(" mb, ");
          Serial.print(P*0.0295333727,2);
          Serial.println(" inHg");

          // The pressure sensor returns abolute pressure, which varies with altitude.
          // To remove the effects of altitude, use the sealevel function and your current altitude.
          // This number is commonly used in weather reports.
          // Parameters: P = absolute pressure in mb, ALTITUDE = current altitude in m.
          // Result: p0 = sea-level compensated pressure in mb

          p0 = pressure.sealevel(P,ALTITUDE); // we're at 1655 meters (Boulder, CO)
          Serial.print("relative (sea-level) pressure: ");
          Serial.print(p0,2);
          Serial.print(" mb, ");
          Serial.print(p0*0.0295333727,2);
          Serial.println(" inHg");
        }
      }
    }
  }
}

void forward_results(int temp, int humidity){
  Serial.print("Temperature = ");
  Serial.println(temp);
  Serial.print("Humidity = ");
  Serial.println(humidity);
}

void print_temp_humidity(){
  Serial.print("th ");
  for(int i=0;i<10;i++){
     if (DHT.read11(DHT11_PIN) == -1){
       forward_results(DHT.temperature,DHT.humidity);
       break;
     }
     delay(10000);
  }
}

void loop(){
  print_all_analog();
  print_temp_preassure();
  print_tempt_humidity();

  delay(60000);
}
