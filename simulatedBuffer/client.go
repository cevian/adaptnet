package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/cevian/adaptnet"
	"github.com/cevian/go-stream/stream"
)

var addr = flag.String("addr", "127.0.0.1:3000", "Server Ip addr")
var targetLevelMs = flag.Int("targetLevelMs", 100, "")
var totalRunMs = flag.Int("totalRunMs", 60000, "")
var probeSize = flag.Bool("probeSize", true, "")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Probe Size", *probeSize)

	runner := stream.NewRunner()

	bf := adaptnet.NewBufferFillProbe(*addr, *probeSize)
	sm := adaptnet.NewSimulatedBuffer(*targetLevelMs, *totalRunMs, bf)
	runner.Add(sm)
	runner.AsyncRunAll()
	runner.Wait()
}
