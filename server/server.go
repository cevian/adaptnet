package main

import (
	"flag"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var numClients = flag.Int("numClients", 2, "")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runnerServer := stream.NewRunner()
	serverOp := adaptnet.NewServerOp(*addr, *numClients)
	runnerServer.Add(serverOp)

	runnerServer.AsyncRunAll()
	runnerServer.Wait()
}
