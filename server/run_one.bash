#!/bin/bash
EXPNAME=$1
PORT=$2
bash run_server_modprobe.bash $EXPNAME $PORT -numClients 1
