#!/bin/bash
EXEC="go run client.go"
SLEEPTIME=20
runint () {
sleep $SLEEPTIME
mkdir data/$1
($EXEC -addr sns58.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $4 -msBetweenChunks $3 2>&1| tee data/$1/client.out; echo "client $(date)" )&
wait
echo "finished"
}

SECS_CHUNK=4
rates=(2350 3600)

for rate in ${rates[*]}
do
  ratebyte=$(($rate*SECS_CHUNK*1000/8))
  echo "rate =" $rate ", ratebyte=" $ratebyte 
  runint $1.$ratebyte.0 $ratebyte 0 10
  runint $1.$ratebyte.1000 $ratebyte 1000 10
done


#235 kbps = 117500 byte chunks
#runint $1.117500.0 117500 0 100
#runint $1.117500.1000 117500 1000 100
#375 kbps = 187500
#runint $1.187500.0 187500 0 100
#runint $1.187500.1000 187500 1000 100
#2350 kbps = 1175000 bytes
#runint $1.1175000.0 1175000 0 10
#runint $1.1175000.1000 1175000 1000 10
#3600 kbps = 1800000 bytes
#runint $1.1800000.0 1800000 0 10
#runint $1.1800000.1000 1800000 1000 10
