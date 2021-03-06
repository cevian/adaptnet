#!/bin/bash
EXEC="go run client.go"
SLEEPTIME=10
runint () {
mkdir data/$1
#(go run client.go -addr sns56.cs.princeton.edu:3001 -bitsPerChunk=100000000 -numChunks $5 -msBetweenChunks 0 2>&1| tee data/$1/client.int.out; echo "interference $(date)" )&
sleep $SLEEPTIME
($EXEC -addr sns58.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $4 -msBetweenChunks $3 2>&1| tee data/$1/client.out; echo "client $(date)" )&
wait
echo "finished"
}

runint $1.10000.0 10000 0 100 1
runint $1.100000.0 100000 0 100 1
runint $1.1000000.0 1000000 0 10 1
runint $1.10000000.0 10000000 0 10 1
runint $1.100000000.0 100000000 0 10 11
runint $1.10000.1000 10000 1000 100 15
runint $1.100000.1000 100000 1000 100 15
runint $1.1000000.1000 1000000 1000 10 15
runint $1.10000000.1000 10000000 1000 10 18
runint $1.100000000.1000 100000000 1000 10 20

