#!/bin/bash

# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../../ && CGO_ENABLED=0 go build -o backend main.go && mv backend ./build/api && cp servidor.env ./build/api && cp -r assets ./build/api && cd ./build/api

# Crea el volúmen de la BD
sudo docker volume create volumen-postgres

# Descomenta línea de DDL de ataque
sed -i '/ataque/s/^#//g' docker-compose.yml

# Crea los contenedores
sudo docker-compose up --detach

# Comenta línea de DDL de ataque
sed -i '/ataque/s/^/#/g' docker-compose.yml

rm ./backend
rm -r ./assets
