/* UDPDaytimeServer
 */
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cevian/adaptnet/udp"
)

func main() {

	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {

	var buf [1400]byte

	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	ch := make(chan []byte, 3)
	go HandleAddrRead(conn, addr, ch)
	HandleAddrWrite(conn, addr, ch)
}

func HandleAddrRead(conn *net.UDPConn, addr *net.UDPAddr, ch chan []byte) {
	defer close(ch)

	for {

		buf := make([]byte, 1440)
		_, _ /*addr*/, err := conn.ReadFromUDP(buf[0:])
		ch <- buf

		if err != nil {
			return
		}
	}
}

func HandleAddrWrite(conn *net.UDPConn, addr *net.UDPAddr, ch chan []byte) {
	for {

		select {
		case buf, ok := <-ch:
			if !ok {
				return
			}
			var report udp.Report
			b := bytes.NewReader(buf)
			err := binary.Read(b, binary.LittleEndian, &report)
			checkError(err)
			fmt.Println("Got report,", report)

		default:
		}
		start := time.Now()
		numPackets := 1

		sendBuf := make([]byte, 1440)
		for i := 0; i < numPackets; i++ {
			dt, err := time.Now().MarshalBinary()
			fmt.Println("Sending, ", i, len(dt))
			checkError(err)
			for index, b := range dt {
				sendBuf[index] = b
			}
			_, err = conn.WriteToUDP(sendBuf, addr)
			if err != nil {
				fmt.Println("Exiting, ", addr)
				return
			}
		}
		runt := time.Now().Sub(start)
		if runt < time.Second {
			time.Sleep(time.Second - runt)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
