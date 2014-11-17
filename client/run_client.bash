#!/bin/bash
EXPNAME=$1
BYTESPERCHUNK=$2
PAUSE=$3
CLIENTARGS="${@:4}"
DATADIR=data/$EXPNAME

mkdir $DATADIR
#sudo tcpdump -i eth0 -w data/$1/tcpdump tcp and port 3000 &
#TCPD=$!

NUMCHUNKS=10
if [ $BYTESPERCHUNK -lt 1000000 ]
then
  NUMCHUNKS=100
fi

go run client.go -addr sns58.cs.princeton.edu:3000 -bitsPerChunk=$BYTESPERCHUNK -numChunks $NUMCHUNKS -msBetweenChunks $PAUSE $CLIENTARGS 2>&1| tee $DATADIR/client.out
#sleep 10
#sudo kill $TCPD
