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
RATE=3.0mbps
# Peak rate to allow
PEAKRATE=3.5mbps
# Interface to shape
IF=eth0
# Average to delay packets by
#LATENCY=100ms
LATENCY=25ms
# Jitter value for packet delay
# Packets will be delayed by $LATENCY +/- $JITTER
JITTER=0ms

#buffer should be > RATE/HZ example  For 10mbit/s on Intel(1000HZ), you need at least 10kbyte buffer if you want to reach your configured rate
#on sns HZ seems to be 125 so for 3MB/s => 3MB/s/125 = 24 kb, double it to be sure
BUFFER=60kb
#MTU as found in ifconfig sns has offloading so you actually want to set this high like 65k
MTU=65000
#Modem q length http://broadband.mpi-sws.org/residential/07_imc_bb.pdf
#amount of time packets can queue before being dropped.
#MODEMQ=60ms
MODEMQ=75ms


start() {
    $TC qdisc add dev $IF root handle 1:0 tbf rate $RATE burst $BUFFER latency $MODEMQ peakrate $PEAKRATE mtu $MTU
    $TC qdisc add dev $IF parent 1:1 handle 10: netem delay $LATENCY $JITTER
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
