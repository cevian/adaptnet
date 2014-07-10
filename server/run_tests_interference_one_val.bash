#!/bin/bash
EXPNAME=$1
PORT=$2
bash run_server_tcpdump.bash $EXPNAME $PORT -numClients 1
