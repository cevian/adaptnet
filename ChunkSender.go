package adaptnet

import (
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type ChunkSender struct {
	client   *netchan.Client
	reader   *netchan.ByteReader
	writer   *netchan.ByteWriter
	response []byte
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

	return &ChunkSender{client, reader, writer, make([]byte, 100)}
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

	t.response, _, err = t.reader.ReadConnectionInto(t.response)
	if err != nil {
		panic(err)
	}

	if len(t.response) != bytesPerChunk {
		panic("Wrong len")
	}

	//tookInternal := time.Since(startInternal)
	return start
}
