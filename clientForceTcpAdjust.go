package adaptnet

import (
	"fmt"
	"sync"
	"time"
)

type ClientForceTcpAdjustOp struct {
	addr                string
	bytesPerChunk       int
	timeBetweenChunksMs int
	numChunksTotal      int
	numChunksOneFlow    int
	maxParallelism      int
}

func NewClientForceTcpAdjustOp(addr string, bytesPerChunk int, timeBetweenChunksMs int, numChunksTotal int, numChunksOneFlow, maxParallelism int) *ClientForceTcpAdjustOp {
	return &ClientForceTcpAdjustOp{addr, bytesPerChunk, timeBetweenChunksMs, numChunksTotal, numChunksOneFlow, maxParallelism}
}

func (t *ClientForceTcpAdjustOp) Run() error {
	css := make([]*ChunkSender, t.maxParallelism)
	for i, _ := range css {
		css[i] = NewChunkSender(t.addr)
	}

	defer func() {
		for _, cs := range css {
			cs.Close()
		}
	}()

	for chunkNo := 0; chunkNo < t.numChunksTotal; chunkNo++ {
		wg := &sync.WaitGroup{}
		wg.Add(len(css))
		start := time.Now()
		cssRun := css
		if chunkNo < t.numChunksOneFlow {
			cssRun = css[:1]
		}
		for _, cs := range cssRun {
			cs_closure := cs
			go func() {
				defer wg.Done()
				cs_closure.MakeRequest(t.bytesPerChunk)
			}()
		}
		wg.Wait()
		took := time.Since(start)

		tookSec := float64(float64(took) / float64(time.Second))
		bandwidthBitsSec := float64(t.bytesPerChunk*len(cssRun)) / tookSec

		fmt.Printf("%d\t%d\t%d\t%E\t%E\t%E\n", len(css), t.timeBetweenChunksMs, t.bytesPerChunk*len(cssRun), float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024))
		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))
	}
	return nil
}

func (t *ClientForceTcpAdjustOp) Stop() error {
	return nil
}
