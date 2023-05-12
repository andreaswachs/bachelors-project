#!/usr/bin/env sh

echo "This setup script will require root access for some commands"

# Installing dependencies
sudo apt update
sudo apt install -y make docker docker-compose

# Configure Docker
sudo systemctl enable docker
sudo systemctl start docker
sudo chmod 666 /var/run/docker.sock
#
# Creating the daaukins user
sudo useradd -m daaukins
usermod -a -G sudo daaukins
usermod -a -G docker daaukins

# ufw should be provided with newer Ubuntu server releases
sudo ufw enable
sudo ufw allow ssh
sudo ufw allow 50052/tcp
sudo ufw allow 40000:50000/tcp

# Checkout the project
git clone https://github.com/andreaswachs/bachelors-project.git /home/daaukins/bachelors-project
ln -s /home/daaukins/bachelors-project/daaukins/server /home/daaukins/server

# Pull docker containers mentioned in yaml files
cd /home/daaukins/server
make pull-images
