package main

import (
	"flag"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var url = flag.String("url", "http://23.239.15.186/app/video/sita/sita.3000_dashinit.mp4", "Url")
var bitsPerChunk = flag.Int("bytesPerChunk", 375000, "")
var msBetweenChunks = flag.Int("msBetweenChunks", 0, "")
var numChunks = flag.Int("numChunks", 100, "")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()

	clientOp := adaptnet.NewClientHttpOp(*url, *bitsPerChunk, *msBetweenChunks, *numChunks)
	runner.Add(clientOp)
	runner.AsyncRunAll()
	runner.Wait()
}
