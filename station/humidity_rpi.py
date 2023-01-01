from datetime import datetime

import serial
#import firebase_admin
#from firebase_admin import credentials
#from firebase_admin import firestore

# from datetime import datetime

#cred = credentials.Certificate("/home/pi/credentials/bewaesserungsanlage-1bb22-firebase-adminsdk-5tuat-e0f60978f7.json")
#firebase_admin.initialize_app(cred)
#db = firestore.client()

import sys
import glob
import serial

# https://stackoverflow.com/questions/12090503/listing-available-com-ports-with-python
def serial_ports():
    """ Lists serial port names

        :raises EnvironmentError:
            On unsupported or unknown platforms
        :returns:
            A list of the serial ports available on the system
    """
    if sys.platform.startswith('win'):
        ports = ['COM%s' % (i + 1) for i in range(256)]
    elif sys.platform.startswith('linux') or sys.platform.startswith('cygwin'):
        # this excludes your current terminal "/dev/tty"
        ports = glob.glob('/dev/tty[A-Za-z]*')
    elif sys.platform.startswith('darwin'):
        ports = glob.glob('/dev/tty.*')
    else:
        raise EnvironmentError('Unsupported platform')

    result = []
    for port in ports:
        try:
            s = serial.Serial(port)
            s.close()
            result.append(port)
        except (OSError, serial.SerialException):
            pass
    return result


def forward_results(serial_port):
    ser = serial.Serial(serial_port, 9600, timeout=1)  # ttyUSB0
    ser.flush()
    line = ser.readline().decode('utf-8').rstrip()

    print(line)
    humidityValues = line.split()
    now = datetime.now()
    date = now.strftime("%Y")+'.'+now.strftime("%m")+'.'+now.strftime("%d")+':'+now.strftime("%H:%M")


if __name__ == '__main__':
    ports=serial_ports()
    print(ports)
    if len(ports)==0:
        forward_results(ports[0])
