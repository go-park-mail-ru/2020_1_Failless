#!/usr/bin/env bash

# Developers accounts
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Серёга Калифорнийский", "password": "qwerty1234", "phone": "88005553535", "email": "almashell@eventum.xyz"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Егор Клинский", "password": "qwerty1234", "phone": "88005553536", "email": "egogoger@eventum.xyz"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Андрей Нижегородский", "password": "qwerty1234", "phone": "88005553537", "email": "rowbotman@eventum.xyz"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Сергей Украинский", "password": "qwerty1234", "phone": "88005553538", "email": "kerch@eventum.xyz"}' \
    http://localhost:3000/api/srv/signup

# Mocks
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Че Гевара", "password": "qwerty1234", "phone": "88000000001", "email": "ernesto@cuba.cu"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Felix Kjellberg", "password": "qwerty1234", "phone": "88000000002", "email": "pewdiepie@sweden.swe"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Владимир Путин", "password": "qwerty1234", "phone": "88000000003", "email": "godhimself@russia.ru"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Евгеша Батиков", "password": "qwerty1234", "phone": "88000000004", "email": "bazhenov@russia.ru"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Линус Торвальдс", "password": "qwerty1234", "phone": "88000000005", "email": "linux@finland.fi"}' \
    http://localhost:3000/api/srv/signup
curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"name": "Григорий Печорин", "password": "qwerty1234", "phone": "88000000006", "email": "thuglife@taman.ta"}' \
    http://localhost:3000/api/srv/signup