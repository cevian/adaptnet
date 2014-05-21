package adaptnet

import (
	"fmt"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type ClientOp struct {
	addr                string
	bytesPerChunk       int
	timeBetweenChunksMs int
	numChunks           int
}

func NewClientOp(addr string, bytesPerChunk int, timeBetweenChunksMs int, numChunks int) *ClientOp {
	return &ClientOp{addr, bytesPerChunk, timeBetweenChunksMs, numChunks}
}

func (t *ClientOp) Run() error {
	client, reader, writer := netchan.NewByteClient(t.addr)
	fmt.Println("Connecting")
	err := client.Connect()
	defer client.Wait()
	defer close(writer)

	if err != nil {
		panic(err)
	}

	for chunkNo := 0; chunkNo < t.numChunks; chunkNo++ {
		r := &Request{int32(t.bytesPerChunk)}
		b, err := SerializeObject(r)
		if err != nil {
			panic(err)
		}
		//fmt.Println("Sending")
		writer <- b
		start := time.Now()

		response := <-reader
		took := time.Since(start)

		if len(response) != t.bytesPerChunk {
			panic("Wrong len")
		}

		tookSec := float64(float64(took) / float64(time.Second))
		bandwidthBitsSec := float64(t.bytesPerChunk) / tookSec

		fmt.Printf("%d\t%E\t%E\t%E\n", t.bytesPerChunk, float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024))
		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))
	}

	fmt.Println("Quiting")
	return nil
}

func (t *ClientOp) Stop() error {
	return nil
}
