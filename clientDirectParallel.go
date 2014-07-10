package adaptnet

import (
	"fmt"
	"sync"
	"time"
)

type ClientDirectParallelOp struct {
	addr                string
	bytesPerChunk       int
	timeBetweenChunksMs int
	numChunks           int
	parallelism         int
}

func NewClientDirectParallelOp(addr string, bytesPerChunk int, timeBetweenChunksMs int, numChunks int, parallelism int) *ClientDirectParallelOp {
	return &ClientDirectParallelOp{addr, bytesPerChunk, timeBetweenChunksMs, numChunks, parallelism}
}

func (t *ClientDirectParallelOp) Run() error {
	css := make([]*ChunkSender, t.parallelism)
	for i, _ := range css {
		css[i] = NewChunkSender(t.addr)
	}

	defer func() {
		for _, cs := range css {
			cs.Close()
		}
	}()

	for chunkNo := 0; chunkNo < t.numChunks; chunkNo++ {
		wg := &sync.WaitGroup{}
		wg.Add(len(css))
		start := time.Now()
		for _, cs := range css {
			cs_closure := cs
			go func() {
				defer wg.Done()
				cs_closure.MakeRequest(t.bytesPerChunk)
			}()
		}
		wg.Wait()
		took := time.Since(start)

		tookSec := float64(float64(took) / float64(time.Second))
		bandwidthBitsSec := float64(t.bytesPerChunk*t.parallelism) / tookSec

		fmt.Printf("%d\t%d\t%d\t%E\t%E\t%E\n", len(css), t.timeBetweenChunksMs, t.bytesPerChunk*t.parallelism, float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024))
		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))
	}
	return nil
}

func (t *ClientDirectParallelOp) Stop() error {
	return nil
}
