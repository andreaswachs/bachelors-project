#!/usr/bin/env sh

block="echo ################################################################################"

echo "This script assumes that it is ran on a newly created ubuntu VM"

# Installing dependencies
$block
echo "Installing dependencies.."
$block
apt update && apt install -y make docker docker-compose

# Configure Docker
$block
echo "Enabling and starting Docker runtime"
$block
systemctl enable docker
systemctl start docker

$block
echo "Changing permissions on docker socket to 666"
$block
chmod 666 /var/run/docker.sock

# Creating the daaukins user
$block
echo "Adding daaukins user and adding daaukins user to docker group"
$block
useradd -m daaukins
usermod -a -G docker daaukins

# ufw should be provided with newer Ubuntu server releases
$block
echo "Enabling ufw firewall and allows daaukins ports. You might need to input 'y' next"
$block
ufw enable
ufw allow ssh
ufw allow 50052/tcp
ufw allow 40000:50000/tcp

# Checkout the project
$block
echo "Downloading Daaukins source code"
$block
git clone https://github.com/andreaswachs/bachelors-project.git /home/daaukins/bachelors-project
ln -s /home/daaukins/bachelors-project/daaukins/server /home/daaukins/server
chown -R daaukins:daaukins /home/daaukins/bachelors-project

# Pull docker containers mentioned in yaml files
$block
echo "Downloading docker images for the Daaukins server"
$block
cd /home/daaukins/server
make pull-images

