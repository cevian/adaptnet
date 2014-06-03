package main

import (
	"flag"
	"math/rand"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var numChunks = flag.Int("numChunks", 10, "")
var numTests = flag.Int("numTests", 100, "")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	runner := stream.NewRunner()

	params := make([]adaptnet.ClientParam, 0)

	sizes := []int{1000000, 100000, 10000}
	pauses := []int{0, 300, 1000}

	for i := 0; i < *numTests; i++ {
		sizei := rand.Intn(len(sizes))
		size := sizes[sizei]

		pausei := rand.Intn(len(pauses))
		pause := pauses[pausei]

		for j := 0; j < *numChunks; j++ {
			params = append(params, adaptnet.ClientParam{size, pause})
		}

	}

	clientOp := adaptnet.NewClientDirectParamOp(*addr, params)
	runner.Add(clientOp)
	runner.AsyncRunAll()
	runner.Wait()
}
