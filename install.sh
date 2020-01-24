#!/bin/sh
set -x
set -e

sudo apt-get update
sudo apt-get upgrade
sudo apt-get install python3-pip
sudo pip3 install adafruit-circuitpython-lis3dh
sudo pip3 install adafruit-circuitpython-ads1x15

make build

sudo cat << EOF > /etc/systemd/system/fruit-pi.service
[Unit]
Description=Fruit-Pi
After=network.target

[Service]
ExecStart=$(pwd)/build/controller/start
WorkingDirectory=$(pwd)/build/controller
StandardOutput=inherit
StandardError=inherit
Restart=always
User=pi

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable fruit-pi.service
sudo systemctl start fruit-pi.service

echo << EOF
The fruit-pi service has been successfully installed. It
will continue to fail unless and until file `~/.env` is
created which declares suitable environment variables. At
present, the required variables are:

FRUIT_PI_HOST
FRUIT_PI_TOKEN
EOF
