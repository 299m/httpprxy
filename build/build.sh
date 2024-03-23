#!/bin/bash

# Add the binary exe to the bin folder
mkdir bin/
cp $(which httpprxy) bin/

# Create the docker image and tag it and push it
sudo docker build -t httpprxy .
sudo docker tag httpprxy:latest 299m/core:httpprxy-1.0.0
sudo docker push 299m/core:httpprxy-1.0.0

# Done