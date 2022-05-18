#!/bin/bash

# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../../ 
CGO_ENABLED=0 go build -o backend main.go 

# Mueve el ejecutable
mv backend ./build/api

# Copia ficheros de variables de entorno
cp envfiles/postgres.env ./build/api
cp envfiles/dns.env ./build/api
cp envfiles/mail.env ./build/api
cp envfiles/servidor.env ./build/api

# Copia los assets
cp -r assets ./build/api

cd ./build/api
# Crea el volúmen de la BD
sudo docker volume create volumen-postgres

# Crea los contenedores
sudo docker-compose up --detach

rm ./backend
rm -r ./assets
rm -r *.env
