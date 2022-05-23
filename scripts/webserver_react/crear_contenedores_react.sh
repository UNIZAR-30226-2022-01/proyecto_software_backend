#!/bin/bash

cd ../../../proyecto_software_frontend_react

npm update
npm run build

cp -r ./build  ../proyecto_software_backend/build/react/web

# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../proyecto_software_backend

CGO_ENABLED=0 go build -o backend main.go

# Mueve el ejecutable
mv backend ./build/react

# Copia ficheros de variables de entorno
cp envfiles/dns.env ./build/react
cp envfiles/servidor.env ./build/react
cp envfiles/clave_tls.key ./build/react
cp envfiles/cert_tls.pem ./build/react

cd ./build/react

sudo docker-compose up --detach

rm ./backend
rm -r ./web/
rm -r *.env
rm -r *.pem
rm -r *.key
