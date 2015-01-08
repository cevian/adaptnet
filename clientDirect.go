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

	fmt.Println("Starting bytesPerChunk", t.bytesPerChunk, "timeBetweenChunksMs", t.timeBetweenChunksMs, "numchunks", t.numChunks)
	err := client.Connect()
	reader := client.Processor().ChannelReader().(*netchan.ByteReader)
	writer := client.Processor().ChannelWriter().(*netchan.ByteWriter)
	defer client.Wait()
	defer writer.CloseConn()

	if err != nil {
		panic(err)
	}

	//response := make([]byte, 100)
	response := make([]byte, 10*1024*1024) //10Mb
        bandwidthBytesSecSum := float64(0)
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

		//response, startInternal, err := reader.ReadConnectionInto(response)
		startInternal, length, err := reader.ReadConnectionIntoNoExpand(response)
		//response, err := reader.ReadConnection()
		if err != nil {
			panic(err)
		}
		took := time.Since(start)
		tookInternal := time.Since(startInternal)

		//if len(response) != t.bytesPerChunk {
		if length != t.bytesPerChunk {
			panic("Wrong len")
		}

		tookSec := float64(float64(took) / float64(time.Second))
		bandwidthBytesSec := float64(t.bytesPerChunk) / tookSec
		bandwidthBytesSecSum += bandwidthBytesSec

		fmt.Printf("%d\t%d\t%d\t%E\t%E\t%E\t%E\t%d\t%E\n", chunkNo, t.timeBetweenChunksMs, t.bytesPerChunk, float64(took), bandwidthBytesSec, bandwidthBytesSec/(1024*1024), float64(tookInternal), start.UnixNano(), 8.0*bandwidthBytesSecSum/float64(chunkNo+1))
		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))
	}
	return nil
}

func (t *ClientDirectOp) Stop() error {
	return nil
}
