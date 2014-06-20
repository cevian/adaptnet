package adaptnet

import (
	"fmt"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type BufferFillProbe struct {
	probeSize      bool
	client         *netchan.Client
	reader         *netchan.ByteReader
	writer         *netchan.ByteWriter
	response       []byte
	rateToUsePerMs int
	baseTimeMs     int

	isProbing    bool
	numProbes    int
	probeBwMsSum float64

	numBase     int
	probingMult int
}

func NewBufferFillProbe(addr string, probeSize bool) *BufferFillProbe {
	client, _, _ := netchan.NewByteClient(addr)
	client.AutoStart = false

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	reader := client.Processor().ChannelReader().(*netchan.ByteReader)
	writer := client.Processor().ChannelWriter().(*netchan.ByteWriter)

	return &BufferFillProbe{probeSize, client, reader, writer, make([]byte, 100), 10, 10, false, 0, 0.0, 0, 10}
}

func (t *BufferFillProbe) Close() {
	t.client.Wait()
	t.writer.CloseConn()
}

func (t *BufferFillProbe) MakeRequest(bytesPerChunk int) time.Duration {
	r := &Request{int32(bytesPerChunk)}
	b, err := SerializeObject(r)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Sending")
	if err := t.writer.WriteConnection(b); err != nil {
		panic(err)
	}
	start := time.Now()

	t.response, _, err = t.reader.ReadConnectionInto(t.response)
	if err != nil {
		panic(err)
	}

	if len(t.response) != bytesPerChunk {
		panic("Wrong len")
	}

	took := time.Since(start)
	//tookInternal := time.Since(startInternal)
	return took
}

func (t *BufferFillProbe) Output(bytesPerChunk int, timePerChunkMs int, took time.Duration) {
	tookSec := float64(float64(took) / float64(time.Second))
	tookMs := float64(float64(took) / float64(time.Millisecond))
	bandwidthBitsSec := float64(bytesPerChunk) / tookSec
	bandwidthBitsMs := float64(bytesPerChunk) / tookMs

	fmt.Printf("%d\t%E\t%E\t%E\t%d\t%d\t%d\n", bytesPerChunk, float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024), timePerChunkMs, t.rateToUsePerMs, int(bandwidthBitsMs))
}

func (t *BufferFillProbe) FillBufferProbing() int {
	timePerChunkMs := t.probingMult * t.baseTimeMs
	t.numProbes++

	bytesPerChunk := t.rateToUsePerMs * timePerChunkMs
	took := t.MakeRequest(bytesPerChunk)
	t.Output(bytesPerChunk, timePerChunkMs, took)

	bufferFillMs := bytesPerChunk / t.rateToUsePerMs

	tookMs := float64(float64(took) / float64(time.Millisecond))
	bandwidthBitsMs := float64(bytesPerChunk) / tookMs
	t.probeBwMsSum += bandwidthBitsMs

	if t.numProbes >= 10 {
		fmt.Println("Debug: Entering base state from probe")
		//return to base state
		bwAvg := float64(t.probeBwMsSum) / float64(t.numProbes)
		if bwAvg > float64(t.rateToUsePerMs)*1.2 {
			fmt.Println("Debug: Changing base state to ", timePerChunkMs)
			t.baseTimeMs = timePerChunkMs
		}
		t.isProbing = false
		t.numBase = 0
	}

	return bufferFillMs
}

func (t *BufferFillProbe) FillBufferBase() int {
	timePerChunkMs := t.baseTimeMs
	t.numBase++

	bytesPerChunk := t.rateToUsePerMs * timePerChunkMs
	took := t.MakeRequest(bytesPerChunk)
	t.Output(bytesPerChunk, timePerChunkMs, took)

	bufferFillMs := bytesPerChunk / t.rateToUsePerMs

	tookMs := float64(float64(took) / float64(time.Millisecond))
	bandwidthBitsMs := float64(bytesPerChunk) / tookMs
	if float64(bandwidthBitsMs) > float64(t.rateToUsePerMs)*1.2 || float64(bandwidthBitsMs) < float64(t.rateToUsePerMs)*0.8 {
		t.rateToUsePerMs = int(bandwidthBitsMs)
		if t.rateToUsePerMs <= 5 {
			t.rateToUsePerMs = 5
		}
	}

	if t.probeSize && t.numBase >= 10 {
		fmt.Println("Debug: Entering probing state from base")
		t.isProbing = true
		t.numProbes = 0
		t.probeBwMsSum = 0
	}

	return bufferFillMs
}

func (t *BufferFillProbe) FillBuffer() int {
	if t.isProbing {
		return t.FillBufferProbing()
	}
	return t.FillBufferBase()
}
