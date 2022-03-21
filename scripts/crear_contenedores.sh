#!/bin/bash

# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../src && CGO_ENABLED=0 go build -o backend main.go && mv backend ../build && cd ../build

# Crea el volúmen de la BD
sudo docker volume create volumen-postgres

# Crea los contenedores
sudo docker-compose up --detach

rm backend
