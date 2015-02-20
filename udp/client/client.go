/* UDPDaytimeClient
 */
package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/cevian/adaptnet/udp"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]

	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	_, err = conn.Write([]byte("anything"))
	checkError(err)

	start := time.Now()
	got := 0
	gotBytes := 0
	totalLatency := time.Duration(0)

	for {

		var buf [1400]byte
		n, err := conn.Read(buf[0:])
		checkError(err)

		got += 1
		gotBytes += n

		var dt time.Time
		err = dt.UnmarshalBinary(buf[0:15])
		checkError(err)
		now := time.Now()
		latency := now.Sub(dt)
		fmt.Println("Latency = ", latency, dt)
		totalLatency += latency

		if now.Sub(start) > time.Second {
			avgLatencyMs := int64(int(totalLatency) / (got * int(time.Millisecond)))
			fmt.Println("Report: avgLatency := ", avgLatencyMs, " got packets: ", got, " bytes ", gotBytes)

			report := &udp.Report{avgLatencyMs, uint32(gotBytes), uint32(got)}
			err := binary.Write(conn, binary.LittleEndian, report)
			checkError(err)

			got = 0
			gotBytes = 0
			start = now
			totalLatency = time.Duration(0)
		}

		//fmt.Println(string(buf[0:n]))

	}

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
