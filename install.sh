#!/bin/bash

sudo apt-get install -y wget xclip 

wget https://github.com/geanify/reaction-dump/releases/download/v0.0.3/reaction-dump
chmod 777 reaction-dump
sudo mv reaction-dump /usr/local/bin/

alias my-reaction-dump='reaction-dump ~/Pictures'
echo "alias my-reaction-dump='reaction-dump ~/Pictures'" >> ~/.bashrc
