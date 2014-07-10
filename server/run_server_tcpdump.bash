#!/bin/bash
# arg1: name of exp
# arg2: port 
# rest passed to script
EXPNAME=$1
PORT=$2
SERVERARGS="${@:3}"
DATADIR=../data/$EXPNAME

mkdir -p $DATADIR
sudo modprobe tcp_probe port=$PORT full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > $DATADIR/server.tcpprobe &
TCPCAP=$!
#sudo tcpdump -i eth0 -w ../data/$1/tcpdump tcp and port 3000 &
#TCPD=$!
go run server.go -addr 0.0.0.0:$PORT $SERVERARGS 2>&1 | tee $DATADIR/server.out
echo "Server pid is $!" 
sudo kill $TCPCAP
#sleep 10
#sudo kill $TCPD
sudo modprobe -r tcp_probe

#analyze tcpdump
#tcpdump -r ../data/$1/tcpdump |less 

#tshark -r aws.v31.10000000.0/tcpdump  -z io,stat,1,"COUNT(tcp.analysis.retransmission) tcp.analysis.retransmission","COUNT(tcp.analysis.duplicate_ack)tcp.analysis.duplicate_ack","COUNT(tcp.analysis.lost_segment) tcp.analysis.lost_segment","COUNT(tcp.analysis.fast_retransmission) tcp.analysis.fast_retransmission"|less
#tshark -r aws.v31.10000000.0/tcpdump -q -z io,stat,1,"COUNT(tcp.analysis.retransmission) tcp.analysis.retransmission","COUNT(tcp.analysis.duplicate_ack)tcp.analysis.duplicate_ack","COUNT(tcp.analysis.lost_segment) tcp.analysis.lost_segment","COUNT(tcp.analysis.fast_retransmission) tcp.analysis.fast_retransmission"|less
