#!/usr/bin/env bash
echo "set -e"
set -e
echo "scp -r ${TRAVIS_BUILD_DIR} a.prokopenko@163.172.133.90:/home/a.prokopenko/eventum/deploy"
scp -i ~/.ssh/deploy_rsa -r $TRAVIS_BUILD_DIR a.prokopenko@163.172.133.90:/home/a.prokopenko/eventum/deploy
echo "ssh a.prokopenko@163.172.133.90 -v exit"
ssh a.prokopenko@163.172.133.90 'cd /home/a.prokopenko/eventum/deploy/2020_1_Failless/deployments;
  cp /home/a.prokopenko/eventum/back/deployments/.env . ;
  sudo docker-compose build ;
  sudo docker-compose up -d ; '
