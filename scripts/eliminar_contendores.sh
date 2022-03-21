#!/bin/bash
cd ../build
sudo docker-compose stop
sudo docker-compose rm
# Borra los volúmenes e imágenes también
sudo docker image rm build_backend postgres
sudo docker volume rm build_volumen-postgres
