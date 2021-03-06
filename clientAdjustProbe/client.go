package main

import (
	"flag"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var msBetweenChunks = flag.Int("msBetweenChunks", 100, "")
var numChunks = flag.Int("numChunks", 100, "")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()

	clientOp := adaptnet.NewClientDirectAdjustProbeOp(*addr, *msBetweenChunks, *numChunks)
	runner.Add(clientOp)
	runner.AsyncRunAll()
	runner.Wait()
}
