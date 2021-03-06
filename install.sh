#!/bin/sh
set -x
set -e

sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install python3-pip -y
sudo pip3 install adafruit-circuitpython-lis3dh
sudo pip3 install adafruit-circuitpython-ads1x15

if [[ ! `go version` = *go1.12.* ]]; then
	wget https://dl.google.com/go/go1.13.7.linux-armv6l.tar.gz
	sudo tar -C /usr/local -xzf go1.13.7.linux-armv6l.tar.gz
fi

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
