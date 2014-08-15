#!/bin/bash
EXEC="go run client.go"
SLEEPTIME=10
runint () {
mkdir data/client.$1
($EXEC -addr sns58.cs.princeton.edu:3000 -bitsPerChunk=$2 -numChunks $4 -msBetweenChunks $3 2>&1| tee data/client.$1/client.out; echo "client $(date)" )&
wait
echo "finished"
}

rates=(250 500 750 1000 1250 1500 1750 2000 5000 10000 15000 20000 30000)
propLat=(50)

for lat in ${propLat[*]}
do
  for rate in ${rates[*]}
  do
    ratebyte=$(($rate*1000))
    echo "rate (KB)" $rate ", ratebyte=" $ratebyte ", latency " $lat
    runint $1.prop.$lat.chunk.$ratebyte.pause.0 $ratebyte 0 10
    runint $1.prop.$lat.chunk.$ratebyte.pause.5000 $ratebyte 5000 10
  done
  mkdir data/client.$1.prop.$lat
  cat data/client.$1.prop.$lat.*/client.out|grep -v "Start" > data/client.$1.prop.$lat/client.out
  sleep $SLEEPTIME
done
