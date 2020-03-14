#!/usr/bin/env bash

sudo apt update -y
sudo apt upgrade -y
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" |sudo tee  /etc/apt/sources.list.d/pgdg.list
sudo apt update -y
sudo apt install -y postgresql-12 postgresql-client-12
systemctl status postgresql.service
sudo apt install -y postgis