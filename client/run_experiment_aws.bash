#!/bin/bash
mkdir data/$1
#sudo tcpdump -i eth0 -w data/$1/tcpdump tcp and port 3000 &
#TCPD=$!

NUMCHUNKS=10
if [ $2 -lt 1000000 ]
then
  NUMCHUNKS=100
fi

go run client.go -addr sns58.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $NUMCHUNKS -msBetweenChunks $3 2>&1| tee data/$1/client.out
#sleep 10
#sudo kill $TCPD
