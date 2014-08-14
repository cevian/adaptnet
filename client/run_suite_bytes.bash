#!/bin/bash
EXEC="go run client.go"
SLEEPTIME=10
runint () {
mkdir data/server.$1
($EXEC -addr sns58.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $4 -msBetweenChunks $3 2>&1| tee data/$1/client.out; echo "client $(date)" )&
wait
echo "finished"
}

rates=(1 5 10 15 20 25 30 35 40 45 50) #mb
propLat=(50)

for lat in ${propLat[*]}
do
  for rate in ${rates[*]}
  do
    ratebyte=$(($rate*1000*1000))
    echo "rate (MB)" $rate ", ratebyte=" $ratebyte ", latency " $lat
    runint $1.prop.$lat.chunk.$ratebyte.pause.0 $ratebyte 0 10
    runint $1.prop.$lat.chunk.$ratebyte.pause.5000 $ratebyte 5000 10
  done
  mkdir data/$1.$lat
  cat data/$1.$lat.*/client.out|grep -v "Start" > data/$1.$lat/client.out
  sleep $SLEEPTIME
done
