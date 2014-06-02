#!/bin/bash
mkdir -p ../data/$1
sudo modprobe tcp_probe port=3000 full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > ../data/$1/client.tcpprobe &
TCPCAP=$!
go run client.go  -addr sns58.cs.princeton.edu:3000 -bitsPerChunk=100000000 -numChunks 10 -msBetweenChunks 0 2>&1 | tee ../data/$1/client.$HOSTNAME.`date +%s`.out
sudo kill $TCPCAP
sudo modprobe -r tcp_probe
