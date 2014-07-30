package adaptnet

import (
	"fmt"
	"sync"
	"time"
)

type ClientDirectParallelAdjustOp struct {
	addr                string
	bytesPerChunk       int
	timeBetweenChunksMs int
	numChunks           int
	targetMs            int
	maxParallelism      int
}

func NewClientDirectParallelAdjustOp(addr string, bytesPerChunk int, timeBetweenChunksMs int, numChunks int, targetMs int, maxParallelism int) *ClientDirectParallelAdjustOp {
	return &ClientDirectParallelAdjustOp{addr, bytesPerChunk, timeBetweenChunksMs, numChunks, targetMs, maxParallelism}
}

func (t *ClientDirectParallelAdjustOp) Run() error {
	css := make([]*ChunkSender, t.maxParallelism)
	for i, _ := range css {
		css[i] = NewChunkSender(t.addr)
	}

	defer func() {
		for _, cs := range css {
			cs.Close()
		}
	}()

	parallelism := 1
	for chunkNo := 0; chunkNo < t.numChunks; chunkNo++ {
		wg := &sync.WaitGroup{}
		wg.Add(parallelism)
		start := time.Now()
		for _, cs := range css[:parallelism] {
			cs_closure := cs
			go func() {
				defer wg.Done()
				cs_closure.MakeRequest(t.bytesPerChunk)
			}()
		}
		wg.Wait()
		took := time.Since(start)

		tookSec := float64(float64(took) / float64(time.Second))
		bandwidthBitsSec := float64(t.bytesPerChunk*parallelism) / tookSec

		fmt.Printf("%d\t%d\t%d\t%E\t%E\t%E\t%E\n", parallelism, t.timeBetweenChunksMs, t.bytesPerChunk*parallelism, float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024), ((float64(t.targetMs)/1000)*bandwidthBitsSec)/float64(t.bytesPerChunk))

		parallelism = int(((float64(t.targetMs) / 1000) * bandwidthBitsSec) / float64(t.bytesPerChunk))
		if parallelism < 1 {
			fmt.Println("DEBUG: parallelism too low ", parallelism)
			parallelism = 1
		}
		if parallelism > t.maxParallelism {
			fmt.Println("DEBUG: parallelism too low", parallelism, t.maxParallelism)
			parallelism = t.maxParallelism
		}

		time.Sleep(time.Millisecond * time.Duration(t.timeBetweenChunksMs))
	}
	return nil
}

func (t *ClientDirectParallelAdjustOp) Stop() error {
	return nil
}
