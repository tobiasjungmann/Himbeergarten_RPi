sudo docker build -t forwarder .
sudo docker rm forwarder .
sudo docker run -it --rm  -e HOME_ASSISTANT_TOKEN=${HOME_ASSISTANT_TOKEN} -p 12348:12348 --name forwarder forwarder
