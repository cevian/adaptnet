package main

import (
	"flag"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var bitsPerChunk = flag.Int("bitsPerChunk", 1000, "")
var msBetweenChunks = flag.Int("msBetweenChunks", 100, "")
var numChunks = flag.Int("numChunks", 100, "")
var direct = flag.Bool("direct", true, "use the direct version of the operator")
var parallelism = flag.Int("par", 0, "# of connections to run in parallel")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()

	if !(*direct) {
		clientOp := adaptnet.NewClientOp(*addr, *bitsPerChunk, *msBetweenChunks, *numChunks)
		runner.Add(clientOp)
	} else {
		if *parallelism == 0 {
			clientOp := adaptnet.NewClientDirectOp(*addr, *bitsPerChunk, *msBetweenChunks, *numChunks)
			runner.Add(clientOp)
		} else {
			clientOp := adaptnet.NewClientDirectParallelOp(*addr, *bitsPerChunk, *msBetweenChunks, *numChunks, *parallelism)
			runner.Add(clientOp)

		}
	}
	runner.AsyncRunAll()
	runner.Wait()
}
