echo "Do not run as sudo, env variable HOME_ASSISTANT_TOKEN is otherwise not set."
sudo docker rm forwarder .
sudo docker build -t forwarder .
sudo docker run -it --restart unless-stopped --rm -e HOME_ASSISTANT_TOKEN=${HOME_ASSISTANT_TOKEN} -p 12348:12348 --name forwarder forwarder
