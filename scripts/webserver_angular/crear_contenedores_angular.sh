#!/bin/bash

cd ../../../proyecto_software_frontend_angular/angular

npm install
npm install @angular/cli
npm update
npm run ng build

cp -r ./dist/angular  ../../proyecto_software_backend/build/angular/web

# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../../proyecto_software_backend

CGO_ENABLED=0 go build -o backend main.go

# Mueve el ejecutable
mv backend ./build/angular

# Copia ficheros de variables de entorno
cp envfiles/dns.env ./build/angular
cp envfiles/servidor.env ./build/angular
cp envfiles/clave_tls.key ./build/angular
cp envfiles/cert_tls.pem ./build/angular

cd ./build/angular

sudo docker-compose up --detach

rm ./backend
rm -r ./web/
rm -r *.env
rm -r *.pem
rm -r *.key
