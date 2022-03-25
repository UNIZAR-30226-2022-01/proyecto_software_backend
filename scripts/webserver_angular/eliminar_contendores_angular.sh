#!/bin/bash

cd ../../build/angular
sudo docker-compose stop
sudo docker-compose rm

# Borra los volúmenes e imágenes también
sudo docker image rm angular_backend

