import serial
import RPi.GPIO as GPIO
from time import sleep

ARDUINO1 = 6
ARDUINO2 = 13



GPIO.setmode(GPIO.BCM)
GPIO.setup(ARDUINO1, GPIO.OUT)
GPIO.setup(ARDUINO2, GPIO.OUT)


if __name__ == '__main__':
        GPIO.output(ARDUINO1, 0)
        GPIO.output(ARDUINO2, 0)        
        GPIO.output(ARDUINO1, 1)
        sleep(4)#ACM0 als ursprungswert
        ser = serial.Serial('/dev/ttyAMA0', 9600, timeout=1)
        ser.flush()
        if ser.in_waiting > 0:
            line = ser.readline().decode('utf-8').rstrip()
            print(line)
        print('finished')
        GPIO.output(ARDUINO1, 0)



