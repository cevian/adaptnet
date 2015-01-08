package adaptnet

import (
	"fmt"
	"net/http"
	"time"
)

type ClientHttpOp struct {
	url                 string
	bytesPerChunk       int
	timeBetweenChunksMs int
	numChunks           int
}

func NewClientHttpOp(url string, bytesPerChunk int, timeBetweenChunksMs int, numChunks int) *ClientHttpOp {
	return &ClientHttpOp{url, bytesPerChunk, timeBetweenChunksMs, numChunks}
}

func (t *ClientHttpOp) Run() error {
	fmt.Println("Starting bytesPerChunk", t.bytesPerChunk, "timeBetweenChunksMs", t.timeBetweenChunksMs, "numchunks", t.numChunks)

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	sumBandwidthBytesSec := float64(0)
	for chunkNo := 0; chunkNo < t.numChunks; chunkNo++ {
		startByte := chunkNo * t.bytesPerChunk
		endByte := (chunkNo + 1) * t.bytesPerChunk

		start := time.Now()

		req, err := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Add("Range", fmt.Sprintf("%d-%d", startByte, endByte))
		resp, err := client.Do(req)
		resp.Body.Close()

		if err != nil {
			panic(err)
		}
		took := time.Since(start)
		tookSec := float64(float64(took) / float64(time.Second))

		bandwidthBytesSec := float64(t.bytesPerChunk) / tookSec
		sumBandwidthBytesSec += bandwidthBytesSec

		fmt.Printf("%d\t%d\t%E\t%E\t%E\t%E\t%d\n", t.timeBetweenChunksMs, t.bytesPerChunk, float64(took), bandwidthBytesSec, bandwidthBytesSec*8, start.UnixNano(), 8*sumBandwidthBytesSec/float64(chunkNo+1))
		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))
	}

	return nil
}

func (t *ClientHttpOp) Stop() error {
	return nil
}
