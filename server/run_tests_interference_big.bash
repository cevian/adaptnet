#!/bin/bash
EXPNAME=$1
PORT=$2
bash run_server_tcpdump.bash $EXPNAME.100000000.0 $PORT -numClients 1
