#!/bin/sh
set -x
set -e

sudo systemctl disable fruit-pi.service
sudo systemctl stop fruit-pi.service

sudo rm /etc/systemd/system/fruit-pi.service
