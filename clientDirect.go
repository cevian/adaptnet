package adaptnet

import (
	"fmt"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type ClientDirectOp struct {
	addr                string
	bytesPerChunk       int
	timeBetweenChunksMs int
	numChunks           int
}

func NewClientDirectOp(addr string, bytesPerChunk int, timeBetweenChunksMs int, numChunks int) *ClientDirectOp {
	return &ClientDirectOp{addr, bytesPerChunk, timeBetweenChunksMs, numChunks}
}

func (t *ClientDirectOp) Run() error {
	client, _, _ := netchan.NewByteClient(t.addr)
	client.AutoStart = false

	fmt.Println("Connecting")
	err := client.Connect()
	reader := client.Processor().ChannelReader().(*netchan.ByteReader)
	writer := client.Processor().ChannelWriter().(*netchan.ByteWriter)
	defer client.Wait()
	defer writer.CloseConn()

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
		if err := writer.WriteConnection(b); err != nil {
			panic(err)
		}
		start := time.Now()

		response, err := reader.ReadConnection()
		if err != nil {
			panic(err)
		}
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

func (t *ClientDirectOp) Stop() error {
	return nil
}
