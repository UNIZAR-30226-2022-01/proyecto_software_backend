#!/bin/bash

cd ../../../proyecto_software_frontend_angular/angular

ng build


cp -r ./dist/angular  ../../proyecto_software_backend/build/angular/web


# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../../proyecto_software_backend/src && CGO_ENABLED=0 go build -o backend main.go && mv backend ../build/angular && cd ../build/angular

sudo docker build -t webserver_angular .
sudo docker run -d -p 8080:8080 --name build_webserver_angular webserver_angular 


rm ./backend
rm -r ./web/
