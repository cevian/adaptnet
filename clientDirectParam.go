package adaptnet

import (
	"fmt"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type ClientParam struct {
	Bytes              int
	TimeBeforeChunksMs int
}

type ClientDirectParamOp struct {
	addr   string
	params []ClientParam
}

func NewClientDirectParamOp(addr string, params []ClientParam) *ClientDirectParamOp {
	return &ClientDirectParamOp{addr, params}
}

func (t *ClientDirectParamOp) Run() error {
	client, _, _ := netchan.NewByteClient(t.addr)
	client.AutoStart = false

	fmt.Println("Starting ClientDirectParamOp with", t.addr)
	err := client.Connect()
	reader := client.Processor().ChannelReader().(*netchan.ByteReader)
	writer := client.Processor().ChannelWriter().(*netchan.ByteWriter)
	defer client.Wait()
	defer writer.CloseConn()

	if err != nil {
		panic(err)
	}

	for _, param := range t.params {
		time.Sleep(time.Millisecond * time.Duration(param.TimeBeforeChunksMs))
		r := &Request{int32(param.Bytes)}
		b, err := SerializeObject(r)
		if err != nil {
			panic(err)
		}
		//fmt.Println("Sending")
		if err := writer.WriteConnection(b); err != nil {
			panic(err)
		}

		start := time.Now()
		response, err := reader.ReadConnection()
		if err != nil {
			panic(err)
		}
		took := time.Since(start)

		if len(response) != param.Bytes {
			panic("Wrong len")
		}

		tookSec := float64(float64(took) / float64(time.Second))
		bandwidthBitsSec := float64(param.Bytes) / tookSec

		fmt.Printf("%d\t%d\t%E\t%E\t%E\n", param.TimeBeforeChunksMs, param.Bytes, float64(took), bandwidthBitsSec, bandwidthBitsSec/(1024*1024))
	}
	return nil
}

func (t *ClientDirectParamOp) Stop() error {
	return nil
}
