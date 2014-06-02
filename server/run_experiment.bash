#!/bin/bash
mkdir -p ../data/$1
sudo modprobe tcp_probe port=3000 full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > ../data/$1/server.tcpprobe &
TCPCAP=$!
go run server.go -numClients 1 -maxPayload 100000000 -addr 0.0.0.0:3000 2>&1 | tee ../data/$1/server.out
sudo kill $TCPCAP
sudo modprobe -r tcp_probe
