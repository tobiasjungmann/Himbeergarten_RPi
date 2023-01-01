const int numberOfAnalogPorts = 8;
void setup() {
  Serial.begin(9600); // open serial port, set the baud rate as 9600 bps
}

void loop() {
  for(int i=0;i<numberOfAnalogPorts;i++){
    Serial.print(analogRead(i));
    Serial.print(" ");
  }

  delay(1000);
}
