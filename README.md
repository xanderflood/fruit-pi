Fruit Pi
--------

A fruiting chamber control module written for a raspberry pi written in go. Uses the following components:

* AM2301 RH/temperature sensor
* Sainsmart 2 channel 5V relay, although most relay modules are probably compatible with this code

Setup
-----

This go binary uses I2C, so you'll need to enable that, and in its current implementation, it also farms out certain tasks to a Python script, meaning that you'll need to install python and verious libraries.

* [This document](https://learn.adafruit.com/circuitpython-on-raspberrypi-linux/installing-circuitpython-on-raspberry-pi) walks you through all things CircuitPython and enabling I2C
** Make sure that `sudo i2cdetect -y 1` shows a device at address 0x48 (assuming you didn't attach the address pin on your ADS1115 to any other pins to change the address).
* Next, you can install the ADS1115 library on top of it following [these instructions](https://github.com/adafruit/Adafruit_CircuitPython_ADS1x15).

(1) Enable I2C in `sudo raspi-config` > Interface options
```
sudo apt-get update
sudo agt-get upgrade
sudo apt-get install python3-pip
sudo pip3 install adafruit-circuitpython-lis3dh
sudo pip3 install adafruit-circuitpython-ads1x15
```
