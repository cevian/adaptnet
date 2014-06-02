#!/bin/bash
mkdir -p ../data/$1
sudo modprobe tcp_probe port=3001 full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > ../data/$1/server.tcpprobe &
TCPCAP=$!
iperf -s -p 3001 | tee ../data/$1/server.iperf.out 
sudo kill $TCPCAP
sudo modprobe -r tcp_probe
