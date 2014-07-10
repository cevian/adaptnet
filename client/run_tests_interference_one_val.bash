#!/bin/bash

runint () {
mkdir data/$1
(go run client.go -addr sns56.cs.princeton.edu:3001 -bitsPerChunk=100000000 -numChunks $5 -msBetweenChunks 0 2>&1| tee data/$1/client.int.out; echo "interference $(date)" )&
sleep 1
(go run client.go -addr sns55.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $4 -msBetweenChunks $3 2>&1| tee data/$1/client.out; echo "client $(date)" )&
wait
echo "finished"
}

runint $1 1000000 0 100 5
