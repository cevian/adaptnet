#!/bin/bash
BITS=100000000
runint () {
mkdir data/$1
(go run client.go -addr sns55.cs.princeton.edu:3001 -bitsPerChunk=$BITS -numChunks $5 -msBetweenChunks 0 2>&1| tee data/$1/client.int.out; echo "interference $(date)" )&
sleep 5
(go run client.go -addr sns56.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $4 -msBetweenChunks $3 2>&1| tee data/$1/client.out; echo "client $(date)" )&
wait
echo "finished"
}

runint $1.$BITS.0 $BITS 0 100 110
