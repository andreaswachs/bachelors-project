#!/usr/bin/env sh

block="################################################################################"
echo="printf $block\n%s\n$block\n\n"

$echo "This script assumes that it is ran on a newly created ubuntu VM"

# Installing dependencies
$echo "Installing dependencies.."
apt update && apt install -y make docker docker-compose

# Configure Docker
$echo "Enabling and starting Docker runtime"
systemctl enable docker
systemctl start docker

$echo "Changing permissions on docker socket to 666"
chmod 666 /var/run/docker.sock

# Creating the daaukins user
$echo "Adding daaukins user and adding daaukins user to docker group"
useradd -m daaukins
usermod -a -G docker daaukins

# ufw should be provided with newer Ubuntu server releases
$echo "Enabling ufw firewall and allows daaukins ports. You might need to input 'y' next"
ufw enable
ufw allow ssh
ufw allow 50052/tcp
ufw allow 40000:50000/tcp

# Checkout the project
$echo "Downloading Daaukins source code"
git clone https://github.com/andreaswachs/bachelors-project.git /home/daaukins/bachelors-project
ln -s /home/daaukins/bachelors-project/daaukins/server /home/daaukins/server
chown -R daaukins:daaukins /home/daaukins/bachelors-project

# Pull docker containers mentioned in yaml files
$echo "Downloading docker images for the Daaukins server"
cd /home/daaukins/server
make pull-images

$echo "Enabling ssh login to daaukins user"
mkdir /home/daaukins/.ssh
chmod 700 /home/daaukins/.ssh
chown daaukins:daaukins /home/daaukins/.ssh
sudo cp /root/.ssh/authorized_keys /home/daaukins/.ssh/authorized_keys
sudo chown -R daaukins:daaukins /home/daaukins/.ssh
sudo chmod 600 /home/daaukins/.ssh/authorized_keys

