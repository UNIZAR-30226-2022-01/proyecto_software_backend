sudo docker rm -f $(sudo docker ps -a -q)
sudo docker volume rm $(sudo docker volume ls -q) 
sudo docker image rm $(sudo docker image ls -q) 
