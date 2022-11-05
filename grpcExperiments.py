#!/home/pi/bin/python
# -*- coding: utf-8 -*-
import os
import sys
from concurrent import futures
import logging
import socket
import re
import time as timetosleep
import argparse

from luma.led_matrix.device import max7219
from luma.core.interface.serial import spi, noop
from luma.core.render import canvas
from luma.core.virtual import viewport
from luma.core.legacy import text, show_message
from luma.core.legacy.font import proportional, CP437_FONT, TINY_FONT, SINCLAIR_FONT, LCD_FONT
from mvg_api import *

import asyncio
import python_weather

from mvg_api import *
from datetime import *

import threading
import ctypes
import socket
import netifaces as ni
import RPi.GPIO as GPIO
import spotipy
from spotipy.oauth2 import SpotifyClientCredentials

import sys

import grpc
from proto.communication_pb2 import GPIOReply, StatusReply
from proto.communication_pb2_grpc import CommunicatorServicer, add_CommunicatorServicer_to_server

# Used GPIO Pins
RELAIS1_GPIO = 19
RELAIS2_GPIO = 26
ARDUINO1_GPIO = 13
ARDUINO2_GPIO = 6

OUTLET_1_GPIO = 14
OUTLET_2_GPIO = 3
OUTLET_3_GPIO = 4
outlets_gpio = [OUTLET_1_GPIO, OUTLET_2_GPIO, OUTLET_3_GPIO]
outlets_state = [False,False,False]

serial = spi(port=0, device=0, gpio=noop())
device = max7219(serial, cascaded=4, block_orientation=-90, rotate=0, blocks_arranged_in_reverse_order=False)
brightness = 255

spotify_scope = 'user-read-currently-playing'
spotify_username = "invalid"
spotify_token = "invalid"  # spotipy.util.prompt_for_user_token(username, scope, "http://127.0.0.1:8080/callback")

PORT_TESTING = 8010
PORT_PRODUCTION = 12345
port = PORT_TESTING

DEFAULT_STARTING_POINT = "Josef Wirt Weg"
DEFAULT_DESTINATION = "Dachau"

matrix_thread = None


class CommunicatorServicer(CommunicatorServicer):
    def outletOn(self, request, context):
        print("Received OutletRequest for ", request.outletId)
        if request.on:
            GPIO.output(outlets_gpio[request.outletId], GPIO.HIGH)
        else:
            GPIO.output(outlets_gpio[request.outletId], GPIO.LOW)
        outlets_state[request.outletId] = request.on
        return GPIOReply(on=outlets_state[request.outletId])
    def getStatus(self, request, context):
        return StatusReply(outlets=outlets_state)


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
# todo problem: innerhalb des files fehlt einmal der richtige import mit.proto


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


def shutdwonCurrentThread(t1):
    if t1 is not None:
        t1.raise_exception()
        t1.join()
        t1 = None


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
    serve(host)
