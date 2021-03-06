#!/bin/bash
usage() {
echo "Usage: shaper.sh {startCombine|startDelay|startRateLimit|stop|show} interface1 interface2"
}



#
#  shaper.sh
#  ---------
#  A utility script for traffic shaping using tc
#
#  Usage
#  -----
#  shape.sh start - starts the shaper
#  shape.sh stop - stops the shaper
#  shape.sh restart - restarts the shaper
#  shape.sh show - shows the rules currently being shaped
#
#  tc uses the following units when passed as a parameter.
#    kbps: Kilobytes per second
#    mbps: Megabytes per second
#    kbit: Kilobits per second
#    mbit: Megabits per second
#    bps: Bytes per second
#  Amounts of data can be specified in:
#    kb or k: Kilobytes
#    mb or m: Megabytes
#    mbit: Megabits
#    kbit: Kilobits
#
#  AUTHORS
#  -------
#  Aaron Blankstein
#  Jeff Terrace
#  Matvey Arye
#
#  Original script written by: Scott Seong
#  Taken from URL: http://www.topwebhosts.org/tools/traffic-control.php
#

# sns notes: 
# 1) TCP (and other) segmentation offloading needs to be turned off.
#    If it isn't then tc sees the huge non-segmented chunks and drops them if needed
#    This cause many tcp packets in a row to be dropped and tcp detects this and goes back into slow-start
#    This is unrealistic.
#
#    To turn off offloading: ethtool -K eth0 tso off; ethtool -K eth0 gro off; ethtool -K eth0 gso off 
#    (not sure that gro and gso need to be disabled but just to be safe)
# 2) Turning on shaping on eth0 in general on sns is dangerous since nfs also uses it. This can cause unexplained slowdown.
#    Thus this script only shapes traffic for one port (PORT)

OFFLOAD=`ethtool -k eth0| grep tcp-segmentation-offload |awk '{print $2}'` 
if [ $OFFLOAD != "off" ]
then
  echo "Offloading is on!. Can't run shaper."
  echo "Please run ethtool -K eth0 tso off; ethtool -K eth0 gro off; ethtool -K eth0 gso off."
  exit
fi



# Name of the traffic control command.
TC=/sbin/tc
# Rate to throttle to
RATE=2.5mbit
# Peak rate to allow
PEAKRATE=3mbit
# Interface to shape
IFS=( "$2" "$3" )
if [ ${#IFS[@]} -lt 1 ] 
then 
usage
exit 1
fi

#IFS=(eth0 eth1)
# Average to delay packets by
#LATENCY=100ms
LATENCY=100ms
# Jitter value for packet delay
# Packets will be delayed by $LATENCY +/- $JITTER
JITTER=2ms

#buffer should be > RATE/HZ example  For 10mbit/s on Intel(1000HZ), you need at least 10kbyte buffer if you want to reach your configured rate
#on sns HZ seems to be 125 so for 3MB/s => 3MB/s/125 = 24 kb, double it to be sure
# for 2.5 mbit/s => 2.5 mbit/s/125 = 20kbit 
BUFFER=60kbit
#MTU as found in ifconfig sns has offloading so you actually want to set this high like 65k
MTU=2000
#Modem q length http://broadband.mpi-sws.org/residential/07_imc_bb.pdf
#amount of time packets can queue before being dropped.
#MODEMQ=60ms
MODEMQ=0ms


startPre(){
    $TC qdisc add dev $IF root handle 1: prio
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3000 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3001 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3002 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3003 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3004 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3005 0xffff flowid 1:1


    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dport 3000 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dport 3001 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dport 3002 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dport 3003 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dport 3004 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dport 3005 0xffff flowid 1:1

    
    #for the socks proxy
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3010 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3011 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3012 0xffff flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3013 0xffff flowid 1:1
    #$TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dst 10.10.0.0/24 flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip dst 192.168.1.107/24 flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip protocol 1 0xFF flowid 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 2 u32 match ip src 0/0 flowid 1:2
}


startCombine() {
    $TC qdisc add dev $IF parent 1:1 handle 10: netem delay 20ms 1ms rate 5mbit limit 10
}

startDelay() {
    $TC qdisc add dev $IF parent 1:1 handle 10: netem delay ${LATENCY}ms ${JITTER}ms limit 100000
}

startDelayOnly() {
    $TC qdisc add dev $IF root handle 1: netem delay 15ms 0ms limit 100000

}
startRateLimit() {
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 4mbit burst 60kbit latency 1000ms peakrate 4.5mbit mtu 2000
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 2mbit burst 60kbit latency 500ms peakrate 2.5mbit mtu 2000
    
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 4mbit burst 60kbit latency 400ms peakrate 4.5mbit mtu 2000
    
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 2mbit burst 60kbit latency 400ms peakrate 4.5mbit mtu 2000
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 2mbit burst 60kbit latency 1000ms peakrate 2.5mbit mtu 2000
    
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 5mbit burst 120kbit latency 800ms peakrate 6mbit mtu 2000
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 3mbit burst 60kbit latency 1000ms peakrate 3.5mbit mtu 2000
    

    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 4mbit burst 60kbit limit 256k peakrate 4.5mbit mtu 2000
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 2mbit burst 60kbit limit 256k peakrate 2.5mbit mtu 2000
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 4mbit burst 60kbit limit 256k  
    #$TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 4mbit burst 120kbit limit 256k  
    $TC qdisc add dev $IF parent 1:1 handle 10: tbf rate 4mbit burst 60kbit limit 256k peakrate 5mbit mtu 1600 
}

startRateLimitNetem() {
    #LATENCY=400
    #RATE=4
    #MTU=1500
    #((LIMIT= (((RATE*1024*1024)/8)*LATENCY/1000)/MTU))
    # echo -n "Starting netem shaping (Latency=$LATENCY RATE=$RATE Limit=$LIMIT): "
    #$TC qdisc add dev $IF parent 1:1 handle 10: netem rate ${RATE}mbit limit $LIMIT
    ((LIMIT= 256000/1500))
    #$TC qdisc add dev $IF parent 1:1 handle 10: netem rate 2mbit limit $LIMIT
    #$TC qdisc add dev $IF parent 1:1 handle 10: netem rate 2mbit limit $LIMIT delay 20ms loss random 0.003%
    #$TC qdisc add dev $IF parent 1:1 handle 10: netem rate 4mbit limit $LIMIT delay 20ms loss random 0.03%
    #$TC qdisc add dev $IF parent 1:1 handle 10: netem rate 4mbit limit $LIMIT delay 20ms loss gemodel 0.3%
    $TC qdisc add dev $IF parent 1:1 handle 10: netem rate 4mbit limit $LIMIT delay 20ms
}


stop() {
    $TC qdisc del dev $IF root
#    $TC qdisc del dev $IF parent 1:1
}


show() {
    $TC -r qdisc ls dev $IF
    echo "--------------------"
    $TC -s qdisc ls dev $IF
}

case "$1" in

startCombine)

echo -n "Starting combined (rate limit and latency) bandwidth shaping: "
for IF in ${IFS[*]} 
do
  startPre
  startCombine
done
echo "done"
;;

startRateLimit)

echo -n "Starting rate limit bandwidth shaping: "
for IF in ${IFS[*]} 
do
  startPre
  startRateLimit
done
echo "done"
;;

startRateLimitNetem)

echo -n "Starting rate limit bandwidth shaping: "
for IF in ${IFS[*]} 
do
  startPre
  startRateLimitNetem
done
echo "done"
;;


startDelay)
LATENCY=${4:-800}
JITTER=${5:-0}
echo -n "Starting latency shaping (Latency=$LATENCY Jitter=$JITTER): "
for IF in ${IFS[*]} 
do
  startPre
  startDelay
done
echo "done"
;;

startDelayOnly)
echo -n "Starting latency shaping only: "
for IF in ${IFS[*]} 
do
  startDelayOnly
done
echo "done"
;;






stop)

echo -n "Stopping bandwidth shaping: "
for IF in ${IFS[*]} 
do
  stop
  done
echo "done"
;;

show)

echo "Bandwidth shaping status for $IF:"
for IF in ${IFS[*]} 
do
  echo "Interface $IF"
  show
  done
echo ""
;;

*)

pwd=$(pwd)
usage
;;

esac 
exit 0
