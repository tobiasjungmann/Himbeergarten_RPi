sudo docker build -t forwarder .
sudo docker rm forwarder .
sudo docker run -it --rm -p 12348:12348 --name forwarder forwarder