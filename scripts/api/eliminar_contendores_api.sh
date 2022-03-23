#!/bin/bash
cd ../../build/api
sudo docker-compose stop
sudo docker-compose rm
# Borra los volúmenes e imágenes también
sudo docker image rm api_backend postgres
sudo docker volume rm api_volumen-postgres volumen-postgres
