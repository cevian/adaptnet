package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var numClients = flag.Int("numClients", 2, "")
var bitsPerChunk = flag.Int("bitsPerChunk", 1000, "")
var msBetweenChunks = flag.Int("msBetweenChunks", 100, "")
var numChunks = flag.Int("numChunks", 100, "")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()
	for i := 0; i < *numClients; i++ {
		clientOp := adaptnet.NewClientOp(*addr, *bitsPerChunk, *msBetweenChunks, *numChunks)
		runner.Add(clientOp)
	}

	runnerServer := stream.NewRunner()
	serverOp := adaptnet.NewServerOp(*addr, *numClients, *bitsPerChunk)
	runnerServer.Add(serverOp)

	runnerServer.AsyncRunAll()
	time.Sleep(time.Second)
	runner.AsyncRunAll()
	runner.Wait()
	runnerServer.Wait()
}
