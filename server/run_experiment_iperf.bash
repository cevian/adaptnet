#!/bin/bash
CONGCONT=${3:-reno}
echo "Using Congestion Control $CONGCONT"
sudo sh -c "echo $CONGCONT > /proc/sys/net/ipv4/tcp_congestion_control"
mkdir -p ../data/$1
sudo modprobe tcp_probe port=$2 full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > ../data/$1/server.tcpprobe &
TCPCAP=$!
( iperf3 -s -p $2 | tee ../data/$1/server.iperf.out ) & 
IPERFD=$!
sleep 10
sudo kill $IPERFD
sudo kill $TCPCAP
sudo modprobe -r tcp_probe
