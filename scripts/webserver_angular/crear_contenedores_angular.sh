#!/bin/bash

cd ../../../proyecto_software_frontend_angular/angular

ng build


cp -r ./dist/angular  ../../proyecto_software_backend/build/angular/web


# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../../proyecto_software_backend && CGO_ENABLED=0 go build -o backend main.go && mv backend ./build/angular && cp servidor.env ./build/angular && cd ./build/angular

sudo docker-compose up --detach

rm ./backend
rm -r ./web/
