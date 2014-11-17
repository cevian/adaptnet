#!/bin/bash
EXPNAME=$1
PORT=${2:-3000}
CONGCONT=${3:-reno}
EXPDIR=../data/$EXPNAME
echo "Using congestion control $CONGCONT"
sudo sh -c "echo $CONGCONT > /proc/sys/net/ipv4/tcp_congestion_control"
sudo sh -c "echo 1 > /proc/sys/net/ipv4/tcp_no_metrics_save"
mkdir $EXPDIR
echo $CONGCONT > $EXPDIR/congcont

rates=(50 100 150 200 250 500 750 1000 1250 1500 1750 2000 5000 10000 20000 30000 40000 50000 60000 70000)
propLat=(50)

for lat in ${propLat[*]}
do
  sudo bash shaper.sh stop eth0 eth1
  sudo bash shaper.sh startDelay eth0 eth1 $lat
  for rate in ${rates[*]}
  do
    ratebyte=$(($rate*1000))
    echo "rate =" $rate ", ratebyte=" $ratebyte ", latency=" $lat
    bash run_server_modprobe.bash $EXPNAME.prop.$lat.chunk.$ratebyte.pause.0 $PORT -numClients 1
    bash run_server_modprobe.bash $EXPNAME.prop.$lat.chunk.$ratebyte.pause.5000 $PORT -numClients 1
  done
done
#cat ../data/$EXPNAME.*/client.out|grep -v "Start" > $EXPDIR/client.out

