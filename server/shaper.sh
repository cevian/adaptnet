#!/bin/bash
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

# Name of the traffic control command.
TC=/sbin/tc
# Rate to throttle to
RATE=3.2mbps
# Peak rate to allow
PEAKRATE=4.0mbps
# Interface to shape
IF=eth0
# Average to delay packets by
#LATENCY=100ms
LATENCY=30ms
# Jitter value for packet delay
# Packets will be delayed by $LATENCY +/- $JITTER
JITTER=2ms

#buffer should be > RATE/HZ example  For 10mbit/s on Intel(1000HZ), you need at least 10kbyte buffer if you want to reach your configured rate
#on sns HZ seems to be 125 so for 3MB/s => 3MB/s/125 = 24 kb, double it to be sure
BUFFER=60kb
#MTU as found in ifconfig sns has offloading so you actually want to set this high like 65k
MTU=2000
#Modem q length http://broadband.mpi-sws.org/residential/07_imc_bb.pdf
#amount of time packets can queue before being dropped.
#MODEMQ=60ms
MODEMQ=60ms


start() {
    #root prio creates classes 1:1, 1:2, and 1:3 (can't create just two bands, 1:3 will be unused)
    $TC qdisc add dev $IF root handle 1: prio
    #port 3000 sent to class 1:1
    $TC filter add dev $IF protocol ip parent 1: prio 1 u32 match ip sport 3000 0xffff flowid 1:1
    #all other traffic sent to class 1:2
    $TC filter add dev $IF protocol ip parent 1: prio 2 u32 match ip src 0/0 flowid 1:2
    
    #attach shapings to class 1:1
    $TC qdisc add dev $IF parent 1:1 handle 10: tbf rate $RATE burst $BUFFER latency $MODEMQ peakrate $PEAKRATE mtu $MTU
    $TC qdisc add dev $IF parent 10:1 handle 101: netem delay $LATENCY $JITTER
#    $TC qdisc add dev $IF root handle 1:0 netem delay $LATENCY $JITTER
#    $TC qdisc add dev $IF parent 1:1 handle 10: tbf rate $RATE burst $BUFFER latency $MODEMQ peakrate $PEAKRATE mtu $MTU
}

stop() {
    $TC qdisc del dev $IF root
    $TC qdisc del dev $IF parent 1:1
}

restart() {
    stop
    sleep 1
    start
}

show() {
    $TC -s qdisc ls dev $IF
}

case "$1" in

start)

echo -n "Starting bandwidth shaping: "
start
echo "done"
;;

stop)

echo -n "Stopping bandwidth shaping: "
stop
echo "done"
;;

restart)

echo -n "Restarting bandwidth shaping: "
restart
echo "done"
;;

show)

echo "Bandwidth shaping status for $IF:"
show
echo ""
;;

*)

pwd=$(pwd)
echo "Usage: shaper.sh {start|stop|restart|show}"
;;

esac 
exit 0
