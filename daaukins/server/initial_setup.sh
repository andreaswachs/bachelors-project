#!/usr/bin/env sh

echo "Ensure you are running this script with root privileges"

apt install -y make docker docker-compose
ufw enable
ufw allow ssh
ufw allow 50052/tcp
ufw allow 40000:50000/tcp
