package main

import (
	"flag"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var numClients = flag.Int("numClients", 2, "")
var maxPayload = flag.Int("maxPayload", 10000000, "")
var direct = flag.Bool("direct", true, "use the direct version of the opertaor")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runnerServer := stream.NewRunner()
	if !(*direct) {
		serverOp := adaptnet.NewServerOp(*addr, *numClients, *maxPayload)
		runnerServer.Add(serverOp)
	} else {
		serverOp := adaptnet.NewServerDirectOp(*addr, *numClients, *maxPayload)
		runnerServer.Add(serverOp)
	}
	runnerServer.AsyncRunAll()
	runnerServer.Wait()
}
