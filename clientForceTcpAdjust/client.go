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
var numChunksTotal = flag.Int("numChunksTotal", 100, "")
var numChunksOneFlow = flag.Int("numChunksOneFlow", 50, "")
var maxParallelism = flag.Int("maxParallelism", 2, "max # of connections to run in parallel")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()

	clientOp := adaptnet.NewClientForceTcpAdjustOp(*addr, *bitsPerChunk, *msBetweenChunks, *numChunksTotal, *numChunksOneFlow, *maxParallelism)
	runner.Add(clientOp)

	runner.AsyncRunAll()
	runner.Wait()
}
