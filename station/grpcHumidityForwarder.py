#!/home/pi/bin/python
# -*- coding: utf-8 -*-
import grpc
import proto.storageServer_pb2
import proto.storageServer_pb2_grpc

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


def forward_results(socket):
    line = socket.readline().decode('utf-8').rstrip()
    print("Line: ",line)
    if len(line)>0:
        humidityValues = line.split()
        print(humidityValues)
        for value in humidityValues:
            response = stub.storeHumidityEntry(proto.storageServer_pb2.StoreHumidityRequest(requestNumber=1, humidity=value))


if __name__ == '__main__':
    print("Started Arduino Analog Pin Forwarder.")
    with grpc.insecure_channel('192.168.178.97:12346') as channel:
        stub = proto.storageServer_pb2_grpc.StorageServerStub(channel)
        ports = serial_ports()
        print(ports)
        serial_instances=[]
        for port in ports:
            serial_instances.append(serial.Serial(port, 9600, timeout=1))
            serial_instances[len(serial_instances)].flush()

        for i in range(4):
            for ser in serial_instances:
                forward_results(ser)
