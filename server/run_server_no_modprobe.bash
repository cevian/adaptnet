#!/bin/bash
EXPNAME=$1
PORT=3000
if [ $# -gt 1  ]
then
  PORT=$2
fi

mkdir -p ../data/$EXPNAME
go run server.go -numClients 1 -maxPayload 100000000 -addr 0.0.0.0:$PORT 2>&1 | tee ../data/$EXPNAME/server.out
