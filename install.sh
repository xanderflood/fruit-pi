#!/bin/sh
set -x
set -e

apt-get update
apt-get upgrade
apt-get install python3-pip
pip3 install adafruit-circuitpython-lis3dh
pip3 install adafruit-circuitpython-ads1x15

make build

cat << EOF > /etc/systemd/system/fruit-pi.service
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

systemctl enable fruit-pi.service
systemctl start fruit-pi.service

echo << EOF
The fruit-pi service has been successfully installed. It
will continue to fail unless and until file `~/.env` is
created which declares suitable environment variables. At
present, the required variables are:

FRUIT_PI_HOST
FRUIT_PI_TOKEN
EOF
