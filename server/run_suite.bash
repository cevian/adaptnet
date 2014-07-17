#!/bin/bash
EXPNAME=$1
PORT=$2
CONGCONT=reno
EXPDIR=../data/$EXPNAME
sudo sh -c "echo $CONGCONT > /proc/sys/net/ipv4/tcp_congestion_control"
mkdir $EXPDIR
cp shaper.sh $EXPDIR/shaper.sh
bash shaper.sh show > $EXPDIR/.shaper.show
echo $CONGCONT > $EXPDIR/congcont
bash run_server_modprobe.bash $EXPNAME.10000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.1000000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.10000000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000000.0 $PORT -numClients 1
#bash run_server_modprobe.bash $EXPNAME.1000000000.0 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.10000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.1000000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.10000000.1000 $PORT -numClients 1
bash run_server_modprobe.bash $EXPNAME.100000000.1000 $PORT -numClients 1
cat ../data/$EXPNAME.*/client.out|grep -v "Start" > $EXPDIR/client.out
#bash run_server_modprobe.bash $EXPNAME.1000000000.1000 $PORT -numClients 1

