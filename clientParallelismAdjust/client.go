package main

import (
	"flag"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var bitsPerChunk = flag.Int("bitsPerChunk", 10000, "")
var msBetweenChunks = flag.Int("msBetweenChunks", 0, "")
var numChunks = flag.Int("numChunks", 100, "")
var targetMs = flag.Int("targetMs", 1000, "")
var maxParallelism = flag.Int("maxParallelism", 100, "# of connections to run in parallel")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()

	clientOp := adaptnet.NewClientDirectParallelAdjustOp(*addr, *bitsPerChunk, *msBetweenChunks, *numChunks, *targetMs, *maxParallelism)
	runner.Add(clientOp)

	runner.AsyncRunAll()
	runner.Wait()
}
