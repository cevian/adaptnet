package adaptnet

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"

	"github.com/cevian/adaptnet/netchan"
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

func (t *ClientDirectAdjustTcpInfoOp) Run() error {
	client, _, _ := netchan.NewByteClient(t.addr)
	client.AutoStart = false

	//fmt.Println("Starting bytesPerChunk", t.bytesPerChunk, "timeBetweenChunksMs", t.timeBetweenChunksMs, "numchunks", t.numChunks)
	err := client.Connect()
	//reader := client.Processor().ChannelReader().(*netchan.ByteReader)
	writer := client.Processor().ChannelWriter().(*netchan.ByteWriter)
	defer client.Wait()
	defer writer.CloseConn()

	if err != nil {
		panic(err)
	}

	f, err := client.NetConn.Conn.(*net.TCPConn).File()
	if err != nil {
		panic(err)
	}
	fd := f.Fd()
	var tcp_info syscall.TCPInfo
	GetTcpInfo(fd, &tcp_info)
	fmt.Printf("tcp_info %+v \n", tcp_info)

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