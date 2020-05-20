#!/usr/bin/env bash
set -e
scp -r $TRAVIS_BUILD_DIR a.prokopenko@163.172.133.90:/home/a.prokopenko/eventum/deploy
ssh a.prokopenko@163.172.133.90 -v exit
cd /home/a.prokopenko/eventum/deploy/2020_1_Failless/deployments
docker-compose build --env-file /home/a.prokopenko/eventum/back/deployments/.env
docker-compose up -d
