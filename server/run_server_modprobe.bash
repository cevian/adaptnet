#!/bin/bash
# arg1: name of exp
# arg2: port 
# rest passed to script
EXPNAME=$1
PORT=${2:-3000}
SERVERARGS="${@:3}"
DATADIR=../data/server.$EXPNAME

mkdir -p $DATADIR
sudo modprobe tcp_probe port=$PORT full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > $DATADIR/server.tcpprobe &
TCPCAP=$!
echo "Server cat pid is $TCPCAP" 
#sudo tcpdump -i eth0 -w ../data/$1/tcpdump tcp and port 3000 &
#TCPD=$!
echo "Starting Server"
./server -addr 0.0.0.0:$PORT $SERVERARGS 2>&1 | tee $DATADIR/server.out
#server will die whel client is done. It is not in background so won't fill $!
sudo kill $TCPCAP
#sleep 10
#sudo kill $TCPD
sudo modprobe -r tcp_probe

#analyze tcpdump
#tcpdump -r ../data/$1/tcpdump |less 

#tshark -r aws.v31.10000000.0/tcpdump  -z io,stat,1,"COUNT(tcp.analysis.retransmission) tcp.analysis.retransmission","COUNT(tcp.analysis.duplicate_ack)tcp.analysis.duplicate_ack","COUNT(tcp.analysis.lost_segment) tcp.analysis.lost_segment","COUNT(tcp.analysis.fast_retransmission) tcp.analysis.fast_retransmission"|less
#tshark -r aws.v31.10000000.0/tcpdump -q -z io,stat,1,"COUNT(tcp.analysis.retransmission) tcp.analysis.retransmission","COUNT(tcp.analysis.duplicate_ack)tcp.analysis.duplicate_ack","COUNT(tcp.analysis.lost_segment) tcp.analysis.lost_segment","COUNT(tcp.analysis.fast_retransmission) tcp.analysis.fast_retransmission"|less
