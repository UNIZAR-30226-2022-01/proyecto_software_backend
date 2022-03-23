#!/bin/bash
sudo docker stop build_webserver_angular
sudo docker rm build_webserver_angular
# Borra los volúmenes e imágenes también
sudo docker image rm webserver_angular
