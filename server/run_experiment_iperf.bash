#!/bin/bash
EXPNAME=$1
PORT=$2
CONGCONT=${3:-reno}
echo "Using Congestion Control $CONGCONT"
sudo sh -c "echo $CONGCONT > /proc/sys/net/ipv4/tcp_congestion_control"
mkdir -p ../data/$EXPNAME
sudo modprobe tcp_probe port=$PORT full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > ../data/$EXPNAME/server.tcpprobe &
TCPCAP=$!
( iperf3 -s -p $PORT | tee ../data/$1/server.iperf.out ) & 
IPERFD=$!
sleep 100
sudo kill $IPERFD
sudo kill $TCPCAP
sudo modprobe -r tcp_probe
