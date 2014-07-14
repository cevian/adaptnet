#!/bin/bash
EXEC="go run client.go"
runint () {
mkdir data/$1
(iperf3 -c sns56.cs.princeton.edu -p 3001 -t 86399 -i 20 -f M -R | tee data/$1/client.iperf.out)&
IPERF=$!
sleep 1
($EXEC -addr sns55.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $4 -msBetweenChunks $3 2>&1| tee data/$1/client.out; echo "client $(date)" )
#kill -- -$IPERF
killall iperf3
wait
echo "finished"
}

runint $1.10000.0 10000 0 100 1
sleep 10
runint $1.100000.0 100000 0 100 1
sleep 10
runint $1.1000000.0 1000000 0 10 1
sleep 10
runint $1.10000000.0 10000000 0 10 1
sleep 10
runint $1.100000000.0 100000000 0 10 11
sleep 10
runint $1.10000.1000 10000 1000 100 15
sleep 10
runint $1.100000.1000 100000 1000 100 15
sleep 10
runint $1.1000000.1000 1000000 1000 10 15
sleep 10
runint $1.10000000.1000 10000000 1000 10 18
sleep 10
runint $1.100000000.1000 100000000 1000 10 20

