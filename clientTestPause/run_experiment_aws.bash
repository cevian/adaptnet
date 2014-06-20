#!/bin/bash
mkdir data/$1
#sudo tcpdump -i eth0 -w data/$1/tcpdump tcp and port 3000 &
#TCPD=$!
go run ./client.go -addr sns58.cs.princeton.edu:3000 -numTests 100 -numChunks 10 2>&1| tee data/$1/client.out
#sleep 10
#sudo kill $TCPD
