#!/bin/bash
mkdir -p ../data/$1
sudo modprobe tcp_probe port=3001 full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > ../data/$1/client.tcpprobe &
TCPCAP=$!
iperf -c sns58.cs.princeton.edu -p 3001 -t 60
sudo kill $TCPCAP
sudo modprobe -r tcp_probe
