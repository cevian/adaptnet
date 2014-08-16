package adaptnet

import (
	"math"
	"sort"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type ChunkSender struct {
	client   *netchan.Client
	reader   *netchan.ByteReader
	writer   *netchan.ByteWriter
	response []byte
	RateLog  []netchan.RateLog
}

func NewChunkSender(addr string) *ChunkSender {
	client, _, _ := netchan.NewByteClient(addr)
	client.AutoStart = false

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	reader := client.Processor().ChannelReader().(*netchan.ByteReader)
	writer := client.Processor().ChannelWriter().(*netchan.ByteWriter)

	return &ChunkSender{client, reader, writer, make([]byte, 100), nil}
}

func (t *ChunkSender) Client() *netchan.Client {
	return t.client
}

func (t *ChunkSender) AvgBandwidth() float64 {
	sum := 0.0
	for _, rle := range t.RateLog {
		rate := float64(rle.Bytes) / (float64(rle.Time) / float64(time.Second))
		sum += rate
	}
	return sum / float64(len(t.RateLog))
}

func (t *ChunkSender) LastBandwidth() float64 {
	rle := t.RateLog[len(t.RateLog)-2]
	rate := float64(rle.Bytes) / (float64(rle.Time) / float64(time.Second))
	return rate
}

func (t *ChunkSender) QuantileBandwidth(quantile int) float64 {
	rates := make([]float64, len(t.RateLog))
	for i, rle := range t.RateLog {
		rate := float64(rle.Bytes) / (float64(rle.Time) / float64(time.Second))
		rates[i] = rate
	}

	sort.Float64s(rates)

	index := (float64(quantile) / 100.0) * float64(len(rates)-1)
	index_ceil := int(math.Ceil(index))
	if len(rates)-1 < index_ceil {
		return rates[len(rates)-1]
	}
	index_floor := int(math.Floor(index))
	return (rates[index_ceil] + rates[index_floor]) / 2.0
}

func (t *ChunkSender) MaxRateLogEntry() *netchan.RateLog {
	var max int
	maxRate := 0.0
	for i, rle := range t.RateLog {
		rate := float64(rle.Bytes) / float64(rle.Time)
		if rate > maxRate {
			max = i
			maxRate = rate
		}
	}
	return &t.RateLog[max]
}

func (t *ChunkSender) Close() {
	t.client.Wait()
	t.writer.CloseConn()
}

func (t *ChunkSender) MakeRequest(bytesPerChunk int) time.Time {
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

	t.response, _, t.RateLog, err = t.reader.ReadConnectionIntoWithLog(t.response, time.Second)

	/*
		for _, rle := range t.RateLog {
			fmt.Println("Bytes", rle.Bytes, "Duration", rle.Time, "Bandwidth (bits/sec)", float64(rle.Bytes*8)/(float64(rle.Time)/float64(time.Second)))
		} */

	if err != nil {
		panic(err)
	}

	if len(t.response) != bytesPerChunk {
		panic("Wrong len")
	}

	//tookInternal := time.Since(startInternal)
	return start
}
