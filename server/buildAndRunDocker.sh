sudo docker build -t plant_server .
sudo docker run -it --rm -p 12346:12346 --privileged --name plant_storage plant_server