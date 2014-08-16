package adaptnet

import (
	"fmt"
	"math"
	"net"
	"syscall"
	"time"
	"unsafe"
)

type ClientDirectAdjustTcpInfoOp struct {
	addr                string
	timeBetweenChunksMs int
	numChunks           int
}

func NewClientDirectAdjustTcpInfoOp(addr string, timeBetweenChunksMs int, numChunks int) *ClientDirectAdjustTcpInfoOp {
	return &ClientDirectAdjustTcpInfoOp{addr, timeBetweenChunksMs, numChunks}
}

func GetTcpInfo(fd uintptr, val *syscall.TCPInfo) (err error) {

	level := syscall.SOL_TCP
	name := syscall.TCP_INFO
	valptr := unsafe.Pointer(val)
	var vallen uint32 = syscall.SizeofTCPInfo

	_, _, e1 := syscall.Syscall6(syscall.SYS_GETSOCKOPT, uintptr(fd), uintptr(level), uintptr(name), uintptr(valptr), uintptr(unsafe.Pointer(&vallen)), 0)
	if e1 != 0 {
		err = e1
	}
	return
}

func NumRttsToBdpAllSS(bdp float64) (rounds float64) {
	startingByte := float64(1500 * 10)
	rnds := math.Log2(bdp / startingByte)
	return math.Ceil(rnds) + 1.0
}

func NumRttsToBdpNoSS(bdp float64) (rounds float64) {
	startingByte := float64(1500 * 10)
	left := bdp - startingByte
	return math.Ceil(left / 1500)
}

func NumRttsToBdp(ssthresh, bdp float64) (rounds float64) {
	toSSthresh := NumRttsToBdpAllSS(ssthresh)
	//startingByte := math.Exp2(toSSthresh - 1.0)*15000.0
	startingByte := ssthresh
	left := bdp - startingByte
	//fmt.Println("ssthresh", toSSthresh, "left", left, left/1500, "ssthresh", ssthresh, "bdp", bdp, math.Exp2(toSSthresh - 1.0)*15000.0, math.Exp2(toSSthresh - 2.0)*15000.0)
	return math.Ceil(left/1500) + toSSthresh
}

func (t *ClientDirectAdjustTcpInfoOp) Run() error {
	cs := NewChunkSender(t.addr)
	defer cs.Close()

	client := cs.Client()

	f, err := client.NetConn.Conn.(*net.TCPConn).File()
	if err != nil {
		panic(err)
	}
	fd := f.Fd()
	var tcp_info syscall.TCPInfo
	GetTcpInfo(fd, &tcp_info)
	//fmt.Printf("tcp_info %+v \n", tcp_info)

	chunkSize := (235 * 1000 * 4) / 8
	for chunkNo := 0; chunkNo < t.numChunks; chunkNo++ {
		start := cs.MakeRequest(chunkSize)

		took := time.Since(start)
		tookSec := float64(float64(took) / float64(time.Second))
		bandwidthBytesSec := float64(chunkSize) / tookSec

		GetTcpInfo(fd, &tcp_info)
		//fmt.Printf("tcp_info %+v \n", tcp_info)
		rtt_us := float64(tcp_info.Rcv_rtt)

		fmt.Printf("%d\t%d\t%E\t%E\t%d\t%E\t", t.timeBetweenChunksMs, chunkSize, float64(took), bandwidthBytesSec, int((bandwidthBytesSec*8)/(1000)), rtt_us/1000)
		//output continues below

		max_rle := cs.MaxRateLogEntry()
		max_bw := float64(max_rle.Bytes) / (float64(max_rle.Time) / float64(time.Second)) //bytes a second
		max_bw = math.Min(cs.QuantileBandwidth(90), bandwidthBytesSec*2)
		max_bw = cs.QuantileBandwidth(90)

		max_bdp := max_bw * (rtt_us / 1000000)

		if max_bw < bandwidthBytesSec {
			fmt.Println("max bw lower than avg bw, shouldnt happen", max_bw, bandwidthBytesSec)
			for _, rle := range cs.RateLog {
				rate := float64(rle.Bytes) / (float64(rle.Time) / float64(time.Second))
				fmt.Println("Entry:", rle.Bytes, rle.Time, rate)
			}
			panic("snh")
		}

		//avg_cwnd := bandwidthBytesSec * rtt_us / 1000000
		//avg_bdp := avg_cwnd
		//minSsthresh := math.Max(avg_cwnd, max_bdp*3/4)
		//last_bdp := cs.LastBandwidth() * rtt_us / 1000000
		//minSsthresh := last_bdp * 3 / 4
		minSsthresh := max_bdp * 3 / 4

		numRoundsToBdp := NumRttsToBdp(minSsthresh, max_bdp)
		numRounds := numRoundsToBdp * 10
		chunkSize = int(numRounds * max_bdp)

		fmt.Printf("%d\t%d\t%d\t%d\t%d\n", int(max_bw*8.0/1000.0), int(max_bdp), int(minSsthresh), int(numRounds), int(chunkSize/1000000))

		//fmt.Println("cs", chunkSize, "max_bw", int(max_bw), "(", int(max_bw*8.0/1000.0), ") max_bdp", max_bdp, "minSsthresh", minSsthresh, "nrtb", numRoundsToBdp, "nr", numRounds)

		//multiplier := 2.0
		/*bdp := bandwidthBytesSec * rtt_us / 1000000

		// BW = cwnd/rtt => rtt * BW = cwnd
		avg_cwnd := bandwidthBytesSec * rtt_us / 1000000
		minSsthresh := avg_cwnd * 3 / 4

		numRoundsToBdp := NumRttsToBdp(minSsthresh, bdp)
		numRoundsToDouble := bandwidthBytesSec/1500


		chunkSize = int(numRoundsToBdp*2+numRoundsToDouble*1.5*bdp)
		fmt.Println("bdp=", bdp, " numRoundsToBdp=", numRoundsToBdp, " NumRoundsToDoube=",numRoundsToDouble, "chunkSize=", chunkSize, " minSsthresh=", minSsthresh)
		*/
		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))

	}

	/*
		//bytesPerChunk := 1000
		response := make([]byte, 100)
		rateToUsePerMs := 1000
		baseTimeMs := 10

		isProbing := false
		numProbes := 0
		probeBwMsSum := 0.0

		numBase := 0
		probingMult := 10
		for chunkNo := 0; chunkNo < t.numChunks; chunkNo++ {
			timePerChunkMs := baseTimeMs
			if isProbing {
				timePerChunkMs = probingMult * baseTimeMs
				numProbes++
			} else {
				numBase++
			}
			bytesPerChunk := rateToUsePerMs * timePerChunkMs

			r := &Request{int32(bytesPerChunk)}
			b, err := SerializeObject(r)
			if err != nil {
				panic(err)
			}
			//fmt.Println("Sending")
			if err := writer.WriteConnection(b); err != nil {
				panic(err)
			}
			start := time.Now()

			response, startInternal, err := reader.ReadConnectionInto(response)
			if err != nil {
				panic(err)
			}
			took := time.Since(start)
			tookInternal := time.Since(startInternal)

			if len(response) != bytesPerChunk {
				panic("Wrong len")
			}

			tookSec := float64(float64(took) / float64(time.Second))
			tookMs := float64(float64(took) / float64(time.Millisecond))
			bandwidthBitsSec := float64(bytesPerChunk) / tookSec
			bandwidthBitsMs := float64(bytesPerChunk) / tookMs

			if isProbing {
				probeBwMsSum += bandwidthBitsMs
			} else {
				if float64(bandwidthBitsMs) > float64(rateToUsePerMs)*1.2 || float64(bandwidthBitsMs) < float64(rateToUsePerMs)*0.8 {
					rateToUsePerMs = int(bandwidthBitsMs)
				}
			}

			fmt.Printf("%d\t%d\t%E\t%E\t%E\t%E\t%d\t%d\t%d\n", t.timeBetweenChunksMs, bytesPerChunk, float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024), float64(tookInternal), timePerChunkMs, rateToUsePerMs, bandwidthBitsMs)


			time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))

			if isProbing && numProbes >= 10 {
				fmt.Println("Debug: Entering base state from probe")
				//return to base state
				bwAvg := float64(probeBwMsSum) / float64(numProbes)
				if bwAvg > float64(rateToUsePerMs)*1.2 {
					fmt.Println("Debug: Changing base state to ", timePerChunkMs)
					baseTimeMs = timePerChunkMs
				}
				isProbing = false
				numBase = 0
			}
			if !isProbing && numBase >= 10 {
				fmt.Println("Debug: Entering probing state from base")
				isProbing = true
				numProbes = 0
				probeBwMsSum = 0
			}
		}*/
	return nil
}

func (t *ClientDirectAdjustTcpInfoOp) Stop() error {
	return nil
}
