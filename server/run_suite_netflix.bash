#!/bin/bash
EXPNAME=$1
PORT=$2
CONGCONT=reno
EXPDIR=../data/$EXPNAME
sudo sh -c "echo $CONGCONT > /proc/sys/net/ipv4/tcp_congestion_control"
mkdir $EXPDIR
cp shaper.sh $EXPDIR/shaper.sh
bash shaper.sh show > $EXPDIR/.shaper.show
echo $CONGCONT > $EXPDIR/congcont

SECONDS=4
rates=(235 375 560 750 1050 1400 1750 2350 3600)

for rate in ${rates[*]}
do
  ratebyte=$((rate*$SECONDS*1000/8))
  bash run_server_modprobe.bash $EXPNAME.$ratebyte.0 $PORT -numClients 1
  bash run_server_modprobe.bash $EXPNAME.$ratebyte.1000 $PORT -numClients 1
done

cat ../data/$EXPNAME.*/client.out|grep -v "Start" > $EXPDIR/client.out

