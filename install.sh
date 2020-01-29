#!/bin/sh
set -x
set -e

sudo apt-get update
sudo apt-get upgrade
sudo apt-get install python3-pip
sudo pip3 install adafruit-circuitpython-lis3dh
sudo pip3 install adafruit-circuitpython-ads1x15

# TODO install go if necessary
/usr/local/go/bin/go build -o build/controller/controller ./cmd/controller/main.go

sudo tee /etc/systemd/system/fruit-pi.service << EOF
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

if sudo systemctl status --no-pager fruit-pi.service; then
	sudo systemctl daemon-reload
	sudo systemctl restart fruit-pi.service
else
	sudo systemctl enable fruit-pi.service
	sudo systemctl start fruit-pi.service
fi

# check whether it's up
sudo systemctl status --no-pager fruit-pi.service

echo << EOF
The fruit-pi service has been successfully installed. It
will continue to fail unless and until file \`~/.env\` is
created which declares suitable environment variables. At
present, the required variables are:

FRUIT_PI_HOST
FRUIT_PI_TOKEN
EOF
