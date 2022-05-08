#!/bin/bash

# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../../ && CGO_ENABLED=0 go build -o backend main.go && mv backend ./build/api && cp servidor.env ./build/api && cp -r assets ./build/api && cd ./build/api

# Crea el volúmen de la BD
sudo docker volume create volumen-postgres

# Crea los contenedores
sudo docker-compose up --detach

rm ./backend
rm -r ./assets
