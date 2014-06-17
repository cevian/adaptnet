package adaptnet

import (
	"fmt"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type ClientDirectAdjustOp struct {
	addr                string
	adjustFreqMs        int
	targetLatencyMs     int
	timeBetweenChunksMs int
	numChunks           int
}

func NewClientDirectAdjustOp(addr string, targetLatencyMs int, timeBetweenChunksMs int, numChunks int) *ClientDirectAdjustOp {
	return &ClientDirectAdjustOp{addr, 0, targetLatencyMs, timeBetweenChunksMs, numChunks}
}

func (t *ClientDirectAdjustOp) Run() error {
	client, _, _ := netchan.NewByteClient(t.addr)
	client.AutoStart = false

	//fmt.Println("Starting bytesPerChunk", t.bytesPerChunk, "timeBetweenChunksMs", t.timeBetweenChunksMs, "numchunks", t.numChunks)
	err := client.Connect()
	reader := client.Processor().ChannelReader().(*netchan.ByteReader)
	writer := client.Processor().ChannelWriter().(*netchan.ByteWriter)
	defer client.Wait()
	defer writer.CloseConn()

	if err != nil {
		panic(err)
	}

	bytesPerChunk := 1000
	for chunkNo := 0; chunkNo < t.numChunks; chunkNo++ {
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

		response, err := reader.ReadConnection()
		if err != nil {
			panic(err)
		}
		took := time.Since(start)

		if len(response) != bytesPerChunk {
			panic("Wrong len")
		}

		tookSec := float64(float64(took) / float64(time.Second))
		tookMs := float64(float64(took) / float64(time.Millisecond))
		bandwidthBitsSec := float64(bytesPerChunk) / tookSec

		fmt.Printf("%d\t%d\t%E\t%E\t%E\n", t.timeBetweenChunksMs, bytesPerChunk, float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024))

		ratio := float64(t.targetLatencyMs) / tookMs
		bytesPerChunk = int(float64(bytesPerChunk) * ratio)

		if bytesPerChunk <= 0 {
			fmt.Printf("Warning Resetting Chunk Sise")
			bytesPerChunkChunk = 10
		}

		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))
	}
	return nil
}

func (t *ClientDirectAdjustOp) Stop() error {
	return nil
}
