#!/bin/bash
sudo docker stop build_webserver_react
sudo docker rm build_webserver_react
# Borra los volúmenes e imágenes también
sudo docker image rm webserver_react
