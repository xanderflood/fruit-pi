#!/bin/sh
set -x
set -e

sudo rm /etc/systemd/system/fruit-pi.service

sudo systemctl disable fruit-pi.service
sudo systemctl stop fruit-pi.service
