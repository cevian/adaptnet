#!/bin/bash
EXPNAME=$1
PORT=$2
sudo sh -c "echo reno > /proc/sys/net/ipv4/tcp_congestion_control"
bash run_server_modprobe.bash $EXPNAME.10000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.1000000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.10000000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.10000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.1000000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.10000000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000000.1000 $PORT -numClients 1
