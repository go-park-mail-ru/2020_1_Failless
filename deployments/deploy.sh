#!/usr/bin/env bash
echo "set -e"
set -e
#echo "scp -r ${TRAVIS_BUILD_DIR} a.prokopenko@163.172.133.90:/home/a.prokopenko/eventum/deploy"
#scp -i ~/.ssh/deploy_rsa -r $TRAVIS_BUILD_DIR a.prokopenko@163.172.133.90:/home/a.prokopenko/eventum/deploy
echo "ssh a.prokopenko@163.172.133.90 -v exit"
ssh a.prokopenko@163.172.133.90 -v exit
echo "cd /home/a.prokopenko/eventum/deploy/2020_1_Failless/deployments"
cd /home/a.prokopenko/eventum/deploy/2020_1_Failless/deployments
echo "docker-compose build --env-file /home/a.prokopenko/eventum/back/deployments/.env"
docker-compose build --env-file /home/a.prokopenko/eventum/back/deployments/.env
echo "docker-compose up -d"
docker-compose up -d
