#!/bin/bash
EXPNAME=$1
PAUSE=$2
CLIENTARGS="${@:3}"
DATADIR=../data/client.$EXPNAME

mkdir $DATADIR
#sudo tcpdump -i eth0 -w data/$1/tcpdump tcp and port 3000 &
#TCPD=$!

NUMCHUNKS=100

git log --pretty=short -10 > $DATADIR/gitversion
git diff > $DATADIR/gitdiff
go run client.go -addr sns58.cs.princeton.edu:3000 -numChunks $NUMCHUNKS -msBetweenChunks $PAUSE $CLIENTARGS 2>&1| tee $DATADIR/client.out
#sleep 10
#sudo kill $TCPD
