package main

import (
	"flag"
	"runtime"
	"fmt"
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

	fmt.Println("Adaptnet", adaptnet.NumRttsToBdp(127833.0, 170444.0))

	clientOp := adaptnet.NewClientDirectAdjustTcpInfoOp(*addr, *msBetweenChunks, *numChunks)
	runner.Add(clientOp)
	runner.AsyncRunAll()
	runner.Wait()
}
