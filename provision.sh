#!/bin/bash
apt-get update -y
apt-get update -y
apt-get upgrade -y
apt-get install -y curl git nano htop autossh

# Install RabbitMQ
wget https://packages.erlang-solutions.com/erlang-solutions_1.0_all.deb
dpkg -i erlang-solutions_1.0_all.deb
wget https://packages.erlang-solutions.com/ubuntu/erlang_solutions.asc
apt-key add erlang_solutions.asc
apt-get update
apt-get install erlang
echo "deb https://dl.bintray.com/rabbitmq/debian xenial main" | sudo tee /etc/apt/sources.list.d/bintray.rabbitmq.list
wget -O- https://www.rabbitmq.com/rabbitmq-release-signing-key.asc | apt-key add -
sudo apt-get update
sudo apt-get install rabbitmq-server


# Install NVM
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.11/install.sh | bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# Install Node & NPM & pm2
nvm install node
npm install -g pm2