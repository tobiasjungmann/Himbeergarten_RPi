# Himbeergarten Plant Utils
(Work in Progress)

This is the backend to the [Himbeergarten App](https://github.com/tobiasjungmann/Himbeergarten_App). It periodically determines the current humidity level of individual plants and stores these values.

The module consists of two parts: the actual server and the measuring station with sensors on the plants.

## Meassuring
Currently, the sensors are connected to Arduinos to use the analog pins. Since they cannot easily be connected to the server directly, 
they are connected serially to a Raspberry Pi, which forwards the information and passes it on to the server. 
To simplify the module, this forwarding will be obsolete by replacing the arduino by an ESP32 with Wifi.

|||
|-|-|
|<img width="800" alt="Standalone module" src="https://github.com/tobiasjungmann/Himbeergarten_RPi/assets/32565407/a4d60792-e4f5-49ff-a694-f93c2238e156">|<img width="800" alt="Module connected" src="https://github.com/tobiasjungmann/Himbeergarten_RPi/assets/32565407/e8b54cf4-5e56-46c6-901e-ac99a7814b03">|


## Storage
Each plant can be stored with multiple images, which can be queried in the original size or in a lower resolution as thumbnails. 
To address the appropriate sensor of a plant, both the corresponding device and the specific GPIO of the plant are stored.

Plants can be added and removed with the controls in the Himbeergarten App.

# Future features
Since this is still a work in progress, I've got a lot of plans:
1. Plant data will be forwarded to an Home Assistant instance.
2. Switch from serially connected Arduinos to ESP32 with Wifi to meassure the humidity directly at the plants.
3. Access token for the API.
