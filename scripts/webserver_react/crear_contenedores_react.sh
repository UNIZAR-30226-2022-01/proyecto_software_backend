#!/bin/bash

cd ../../../proyecto_software_frontend_react

npm run build



cp -r ./build  ../proyecto_software_backend/build/react/web


# Compila el backend sin depender de librerías de C y trae el ejecutable a la carpeta local
cd ../proyecto_software_backend/src && CGO_ENABLED=0 go build -o backend main.go && mv backend ../build/react && cd ../build/react

#sudo docker build -t webserver_react .
#sudo docker run -d -p 8080:8080 --name build_webserver_react webserver_react

sudo docker-compose up --detach

rm ./backend
rm -r ./web/
