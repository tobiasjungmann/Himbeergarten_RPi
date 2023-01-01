void setup() {
  Serial.begin(9600); // open serial port, set the baud rate as 9600 bps
}

void loop() {
  int val0 = analogRead(0); //connect sensor to Analog 0
  int val1 = analogRead(1); //connect sensor to Analog 0
  int val2 = analogRead(2); //connect sensor to Analog 0
  int val3 = analogRead(3); //connect sensor to Analog 0
  int val4 = analogRead(4); //connect sensor to Analog 0
  Serial.print("Meassured: "); //print the value to serial port
  Serial.print(val0);
  Serial.print(" ");
  Serial.print(val1);
  Serial.print(" ");
  Serial.print(val2);
  Serial.print(" ");
  Serial.print(val3);
  Serial.print(" ");
  Serial.print(val4);
  Serial.println("Neue Zeile");
  delay(500);
}
