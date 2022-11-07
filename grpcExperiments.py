#!/home/pi/bin/python
# -*- coding: utf-8 -*-
from concurrent import futures

from luma.led_matrix.device import max7219
from luma.core.interface.serial import spi, noop

import socket
import RPi.GPIO as GPIO

import grpc

from matrix.showMVV import showMVV
from matrix.showSpotify import showSpotify
from matrix.showWeather import show_weather
from proto.communication_pb2 import GPIOReply, StatusReply, EmptyMsg, MatrixState, MatrixChangeModeReply
from proto.communication_pb2_grpc import CommunicatorServicer, add_CommunicatorServicer_to_server
from matrix.showTime import showTime

# Used GPIO Pins
RELAIS1_GPIO = 19
RELAIS2_GPIO = 26

ARDUINO1_GPIO = 13
ARDUINO2_GPIO = 6
OUTLET_1_GPIO = 14
OUTLET_2_GPIO = 3
OUTLET_3_GPIO = 4
outlets_gpio = [OUTLET_1_GPIO, OUTLET_2_GPIO, OUTLET_3_GPIO, ARDUINO1_GPIO, ARDUINO2_GPIO]
outlets_state = [False, False, False, False, False]
matrix_thread_array = [None]
matrix_state = [MatrixState.MATRIX_TIME]

serial = spi(port=0, device=0, gpio=noop())
device = max7219(serial, cascaded=4, block_orientation=-90, rotate=0, blocks_arranged_in_reverse_order=False)
brightness = 255

PORT_TESTING = 8010
PORT_PRODUCTION = 12345
port = PORT_TESTING


class CommunicatorServicer(CommunicatorServicer):
    def outletOn(self, request, context):
        print("Received OutletRequest for ", request.id)
        if request.on:
            GPIO.output(outlets_gpio[request.id], GPIO.HIGH)
        else:
            GPIO.output(outlets_gpio[request.id], GPIO.LOW)
        outlets_state[request.id] = request.on
        return GPIOReply(statusList=outlets_state)

    def getStatus(self, request, context):
        return StatusReply(gpios=outlets_state, matrixState=matrix_state[0],matrixBrightness=brightness)

    def matrixSetActivated(self, request, context):
        return EmptyMsg()

    def matrixSetMode(self, request, context):
        if matrix_thread_array[0] is not None:
            matrix_thread_array[0].raise_exception()
            matrix_thread_array[0].join()
            matrix_thread_array[0] = None
        matrix_state[0]=request.state
        if request.state == MatrixState.MATRIX_TIME:
            matrix_thread_array[0] = showTime('Thread 1', device)
        elif request.state == MatrixState.MATRIX_WEATHER:
            matrix_thread_array[0] = show_weather('Thread 1', device)
        elif request.state == MatrixState.MATRIX_SPOTIFY:
            matrix_thread_array[0] = showSpotify('Thread 1', device)
        elif request.state == MatrixState.MATRIX_MVV:
            matrix_thread_array[0] = showMVV('Thread 1', device, request.start, request.destination)
        elif request.state == MatrixState.MATRIX_QUIT:
            exit()

        if request.state == MatrixState.MATRIX_STANDBY:
            device.cleanup()
        else:
            matrix_thread_array[0].start()

        return MatrixChangeModeReply(state=request.state)

    def matrixSetBrightness(self, request, context):
        return EmptyMsg()


def serve(address):
    print("Starting Server...")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_CommunicatorServicer_to_server(
        CommunicatorServicer(), server)
    concataddress = address + ":" + str(port)
    print("Server at: ", concataddress)
    server.add_insecure_port(concataddress)
    server.start()
    print("Waiting for connections...")
    server.wait_for_termination()
    print("Terminated.")


def get_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.settimeout(0)
    try:
        # doesn't even have to be reachable
        s.connect(('10.255.255.255', 1))
        IP = s.getsockname()[0]
    except Exception:
        IP = '127.0.0.1'
    finally:
        s.close()
    return IP


def init_gpios():
    GPIO.setmode(GPIO.BCM)

    GPIO.setup(RELAIS1_GPIO, GPIO.OUT)
    GPIO.setup(RELAIS2_GPIO, GPIO.OUT)
    GPIO.setup(ARDUINO1_GPIO, GPIO.OUT)
    GPIO.setup(ARDUINO2_GPIO, GPIO.OUT)

    for outlet in outlets_gpio:
        GPIO.setup(outlet, GPIO.OUT)
        GPIO.output(outlet, GPIO.LOW)

    GPIO.output(RELAIS1_GPIO, GPIO.LOW)
    GPIO.output(RELAIS2_GPIO, GPIO.LOW)
    GPIO.output(ARDUINO1_GPIO, GPIO.LOW)
    GPIO.output(ARDUINO2_GPIO, GPIO.LOW)


if __name__ == '__main__':
    host = get_ip()
    print(host)
    init_gpios()
    matrix_thread_array[0] = showTime('Thread 1', device)
    matrix_thread_array[0].start()
    serve(host)
