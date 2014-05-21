package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()
	runnerServer := stream.NewRunner()
	clientOp := adaptnet.NewClientOp(*addr)
	serverOp := adaptnet.NewServerOp(*addr)
	runner.Add(clientOp)
	runnerServer.Add(serverOp)

	runnerServer.AsyncRunAll()
	time.Sleep(time.Second)
	runner.AsyncRunAll()
	runner.Wait()
	runnerServer.Wait()
}
