#!/bin/bash
mkdir -p ../data/$1
sudo modprobe tcp_probe port=3000 full=1
sudo chmod 444 /proc/net/tcpprobe
cat /proc/net/tcpprobe > ../data/$1/server.tcpprobe &
TCPCAP=$!
sudo tcpdump -i eth0 -w ../data/$1/tcpdump tcp and port 3000 &
TCPD=$!
go run server.go -numClients 1 -maxPayload 100000000 -addr 0.0.0.0:3000 2>&1 | tee ../data/$1/server.out
echo "Server pid is $!" 
sudo kill $TCPCAP
sleep 10
sudo kill $TCPD
sudo modprobe -r tcp_probe

#analyze tcpdump
#tcpdump -r ../data/$1/tcpdump |less 

#tshark -r aws.v31.10000000.0/tcpdump  -z io,stat,1,"COUNT(tcp.analysis.retransmission) tcp.analysis.retransmission","COUNT(tcp.analysis.duplicate_ack)tcp.analysis.duplicate_ack","COUNT(tcp.analysis.lost_segment) tcp.analysis.lost_segment","COUNT(tcp.analysis.fast_retransmission) tcp.analysis.fast_retransmission"|less
#tshark -r aws.v31.10000000.0/tcpdump -q -z io,stat,1,"COUNT(tcp.analysis.retransmission) tcp.analysis.retransmission","COUNT(tcp.analysis.duplicate_ack)tcp.analysis.duplicate_ack","COUNT(tcp.analysis.lost_segment) tcp.analysis.lost_segment","COUNT(tcp.analysis.fast_retransmission) tcp.analysis.fast_retransmission"|less
