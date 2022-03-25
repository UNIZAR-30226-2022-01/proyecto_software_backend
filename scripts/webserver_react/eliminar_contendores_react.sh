#!/bin/bash

cd ../../build/react
sudo docker-compose stop
sudo docker-compose rm

# Borra los volúmenes e imágenes también
sudo docker image rm react_backend
