# Himbeergarten RPi

## Humidity Utils
This module controls the humidity level of indoor plants and writes these values in regular intervals into a Firebase database.
Therefore, multiple Arduinos can be connected to a central RaspberryPi.

## Matrix Utils
Displays time, currently played songs and departure information on a MAX7219 LED Matrix.
Additionally, relais can be connected to control external outlets as well as lights. 

## Setup and Usage
Setup "spotipy":
pip install spotipy
export SPOTIPY_CLIENT_ID='your-spotify-client-id'
export SPOTIPY_CLIENT_SECRET='your-spotify-client-secret'
export SPOTIPY_REDIRECT_URI='your-app-redirect-url'
Initialize spotify credentials redirect uri. Afterwards fill in the credentials in the browser window

				
Setup MAX7219 Library: Follow the steps in the [installation manual](https://luma-led-matrix.readthedocs.io/en/latest/install.html) and connect the matrix to the correct GPIO Pins


install "MVV API": pip install mvg-api
install "netifaces": pip install netifaces

Start with: python3 rpiReceiver.py
