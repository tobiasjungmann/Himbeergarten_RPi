#!/home/pi/bin/python
# -*- coding: utf-8 -*-

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

# Used GPIO Pins
RELAIS1 = 19
RELAIS2 = 26
ARDUINO1 = 13
ARDUINO2 = 6
OUTLET1 = 14
OUTLET2 = 3
OUTLET3 = 4

GPIO.setmode(GPIO.BCM)

GPIO.setup(RELAIS1, GPIO.OUT)
GPIO.setup(RELAIS2, GPIO.OUT)
GPIO.setup(ARDUINO1, GPIO.OUT)
GPIO.setup(ARDUINO2, GPIO.OUT)
GPIO.setup(OUTLET1, GPIO.OUT)
GPIO.setup(OUTLET2, GPIO.OUT)
GPIO.setup(OUTLET3, GPIO.OUT)

GPIO.output(RELAIS1, GPIO.LOW)
GPIO.output(RELAIS2, GPIO.LOW)
GPIO.output(ARDUINO1, GPIO.LOW)
GPIO.output(ARDUINO2, GPIO.LOW)
GPIO.output(OUTLET1, GPIO.LOW)
GPIO.output(OUTLET2, GPIO.LOW)
GPIO.output(OUTLET3, GPIO.LOW)

relais1activated = False
relais2activated = False
arduino1activated = False
arduino2activated = False
outlet1activated = False
outlet2activated = False
outlet3activated = False

serial = spi(port=0, device=0, gpio=noop())
device = max7219(serial, cascaded=4, block_orientation=-90, rotate=0, blocks_arranged_in_reverse_order=False)
brightness = 255

scope = 'user-read-currently-playing'
username = "invalid"
token = "invalid"  # spotipy.util.prompt_for_user_token(username, scope, "http://127.0.0.1:8080/callback")

port = 15439

startingPoint = "Josef Wirt Weg"
destination = "Dachau"

t1 = None


# Retrieves the current IP address
# Source: https://stackoverflow.com/questions/166506/finding-local-ip-addresses-using-pythons-stdlib
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


# Loosly follows https://github.com/vierofernando/python-weather
async def getweather(device):
    client = python_weather.Client(format=python_weather.METRIC)
    weather = await client.find("Munich")
    weatherString = str(weather.current.temperature) + "°C " + str(weather.current.feels_like) + " " + str(
        weather.current.wind_display)
    show_message(device, weatherString, fill="white", font=proportional(CP437_FONT))#, scroll_delay=0.02)


# show song title between time
class show_weather(threading.Thread):
    def __init__(self, name):
        threading.Thread.__init__(self)
        self.name = name

    def run(self):
        try:
            while True:
                loopa = asyncio.new_event_loop()
                loopa.run_until_complete(getweather(device))
                timetosleep.sleep(2)
        finally:
            print('ended')

    def get_id(self):

        # returns id of the respective thread
        if hasattr(self, '_thread_id'):
            return self._thread_id
        for id, thread in threading._active.items():
            if thread is self:
                return id

    def raise_exception(self):
        thread_id = self.get_id()
        res = ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id,
                                                         ctypes.py_object(SystemExit))
        if res > 1:
            ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, 0)
            print('Exception raise failure')


# shows the time while running
class show_route(threading.Thread):
    def __init__(self, name):
        threading.Thread.__init__(self)
        self.name = name

    def run(self):

        # target function of the thread class
        try:
            while True:
                while True:
                    routes = get_route(get_id_for_station(startingPoint), get_id_for_station(destination), time=None,
                                       arrival_time=False, max_walk_time_to_start=None, max_walk_time_to_dest=None,
                                       change_limit=None, ubahn=True, bus=True, tram=True, sbahn=True)

                    show_message(device, startingPoint + "  ->  " + destination, fill="white",
                                 font=proportional(CP437_FONT))

                    for route in routes:
                        departure=route['departure_datetime']
                        arrival=route['arrival_datetime']
                        with canvas(device) as draw:
                            text(draw, (0, 0), str(departure.hour).zfill(2)+":"+str(departure.minute).zfill(2), fill="white", font=proportional(CP437_FONT))
                        timetosleep.sleep(3)
                        with canvas(device) as draw:
                            device.contrast(50)
                            text(draw, (0, 0), str(arrival.hour).zfill(2)+":"+str(arrival.minute).zfill(2), fill="white", font=proportional(CP437_FONT))
                        timetosleep.sleep(1)
                        with canvas(device) as draw:
                            device.contrast(200)
        finally:
            print('ended')


    def get_id(self):
       # returns id of the respective thread
        if hasattr(self, '_thread_id'):
            return self._thread_id
        for id, thread in threading._active.items():
            if thread is self:
                return id


    def raise_exception(self):
        thread_id = self.get_id()
        res = ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id,
                                                     ctypes.py_object(SystemExit))
        if res > 1:
            ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, 0)
            print('Exception raise failure')


# shows the time while running
class showTime(threading.Thread):
    def __init__(self, name):
        threading.Thread.__init__(self)
        self.name = name

    def run(self):
        try:
            while True:
                t = datetime.now(tz=None)
                with canvas(device) as draw:
                    text(draw, (0, 0), str(t.hour).zfill(2) + ":" + str(t.minute).zfill(2), fill="white",
                         font=proportional(CP437_FONT))
                timetosleep.sleep(1)
        finally:
            print('ended')

    def get_id(self):

        # returns id of the respective thread
        if hasattr(self, '_thread_id'):
            return self._thread_id
        for id, thread in threading._active.items():
            if thread is self:
                return id

    def raise_exception(self):
        thread_id = self.get_id()
        res = ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id,
                                                         ctypes.py_object(SystemExit))
        if res > 1:
            ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, 0)
            print('Exception raise failure')


# show song title between time
class show_songTitle(threading.Thread):
    def __init__(self, name):
        threading.Thread.__init__(self)
        self.name = name

    def run(self):
        try:
            while True:
                sp = spotipy.Spotify(auth=token)
                current_song = sp.current_user_playing_track()
                if current_song != None:
                    print(current_song['item']['name'])
                    show_message(device, current_song['item']['name'], fill="white", font=proportional(CP437_FONT))
                for i in range(0, 20):
                    t = datetime.now(tz=None)
                    text(device, (0, 0), str(t.hour).zfill(2) + ":" + str(t.minute).zfill(2), fill="white",
                        font=proportional(CP437_FONT))
                    timetosleep.sleep(1)
        finally:
            print('ended')

    def get_id(self):

        # returns id of the respective thread
        if hasattr(self, '_thread_id'):
            return self._thread_id
        for id, thread in threading._active.items():
            if thread is self:
                return id

    def raise_exception(self):
        thread_id = self.get_id()
        res = ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id,
                                                         ctypes.py_object(SystemExit))
        if res > 1:
            ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, 0)
            print('Exception raise failure')

def shutdwonCurrentThread(t1):
    if t1 is not None:
        t1.raise_exception()
        t1.join()
        t1 = None

soc = socket.socket()
host = get_ip()

print("Starting to show IP")
show_message(device, "IP " + host, fill="white", font=proportional(CP437_FONT))#, scroll_delay=0.1)
print("IP message ended")

soc.bind((host, port))
soc.listen(5)
conn, addr = soc.accept()


try:
    while True:
        print("Got connection from: ", addr)
        length_of_message = int.from_bytes(conn.recv(2), byteorder='big')
        msg = conn.recv(length_of_message).decode("UTF-8")
        print(msg)

        if "Stations" in msg:
            startingPoint = msg.split(";")[1]
            destination = msg.split(";")[2]
            message_to_send = "Stations were received".encode("UTF-8")
            shutdwonCurrentThread(t1)
            t1 = show_route('Thread 1')
            t1.start()

        elif msg == "changetime":
            message_to_send = "changeTime was received".encode("UTF-8")
            shutdwonCurrentThread(t1)
            t1 = showTime('Thread 1')
            t1.start()
        elif msg == "songtitle":
            message_to_send = "spotipy was received".encode("UTF-8")
            shutdwonCurrentThread(t1)
            t1 = show_songTitle('Thread 1')
            t1.start()
        elif msg == "weather":
            message_to_send = "weather was received".encode("UTF-8")
            shutdwonCurrentThread(t1)
            t1 = show_weather('Thread 1')
            t1.start()


        elif "newBrightness" in msg:
            message_to_send = "brightness was received".encode("UTF-8")
            brightness = (int)(msg.split(":")[1])
            print(brightness)
            device.contrast(brightness)

        elif msg == "relais1":
            message_to_send = "relais1 was received".encode("UTF-8")
            if relais1activated:
                GPIO.output(RELAIS1, GPIO.LOW)
            else:
                GPIO.output(RELAIS1, GPIO.HIGH)
            relais1activated = not relais1activated

        elif msg == "relais2":
            message_to_send = "relais2 was received".encode("UTF-8")
            if relais2activated:
                GPIO.output(RELAIS2, GPIO.LOW)
            else:
                GPIO.output(RELAIS2, GPIO.HIGH)
            relais2activated = not relais2activated
        elif msg == "arduino1":
            message_to_send = "arduino1 was received".encode("UTF-8")
            if arduino1activated:
                GPIO.output(ARDUINO1, GPIO.LOW)
            else:
                GPIO.output(ARDUINO1, GPIO.HIGH)
            arduino1activated = not arduino1activated
        elif msg == "arduino2":
            message_to_send = "arduino2 was received".encode("UTF-8")
            if arduino2activated:
                GPIO.output(ARDUINO2, GPIO.LOW)
            else:
                GPIO.output(ARDUINO2, GPIO.HIGH)
            arduino2activated = not arduino2activated

        elif msg == "outlet1":
            message_to_send = "outlet1 was received".encode("UTF-8")
            if outlet1activated:
                GPIO.output(OUTLET1, GPIO.LOW)
            else:
                GPIO.output(OUTLET1, GPIO.HIGH)
            outlet1activated = not outlet1activated
        elif msg == "outlet2":
            message_to_send = "outlet2 was received".encode("UTF-8")
            if outlet2activated:
                GPIO.output(OUTLET2, GPIO.LOW)
            else:
                GPIO.output(OUTLET2, GPIO.HIGH)
            outlet2activated = not outlet2activated
        elif msg == "outlet3":
            message_to_send = "outlet3 was received".encode("UTF-8")
            if outlet3activated:
                GPIO.output(OUTLET3, GPIO.LOW)
            else:
                GPIO.output(OUTLET3, GPIO.HIGH)
            outlet3activated = not outlet3activated


        elif msg == "standby":
            message_to_send = "deaktivate was received".encode("UTF-8")
            shutdwonCurrentThread(t1)
            device.cleanup()
        elif msg == "quit":
            message_to_send = "shutdown was received".encode("UTF-8")
            soc.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            soc.shutdown(1)
            exit()
        else:
            message_to_send = "Error: something wrong was send.".encode("UTF-8")

        conn, addr = soc.accept()
        conn.send(len(message_to_send).to_bytes(2, byteorder='big'))
        conn.send(message_to_send)
finally:
    print("             terminated")
    soc.close()
