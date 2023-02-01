#!/bin/bash

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

wget "https://github.com/Kitchen-Kreations/listparse/releases/download/v1.0.0/listparse" -O /usr/sbin