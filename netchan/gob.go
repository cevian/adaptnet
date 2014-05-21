package netchan

import (
	"encoding/gob"
	"fmt"
	"io"
	"sync"

	"github.com/cevian/go-stream/stream"
)

type GobProcessorGenerator struct{}

func (t GobProcessorGenerator) NewChannelProcessor() ChannelProcessor {
	return &SimpleChannelProcessor{&sync.WaitGroup{}, NewGobReader(), NewGobWriter()}
}

func NewGobClient(addr string) *Client {
	return NewClient(addr, GobProcessorGenerator{}.NewChannelProcessor())
}

func NewGobServer(addr string) *Server {
	return NewServer(addr, GobProcessorGenerator{})
}

type GobWriter struct {
	channel chan stream.Object
	conn    Conn
}

func NewGobWriter() *GobWriter {
	return &GobWriter{make(chan stream.Object, 3), nil}
}

func (t *GobWriter) SetConn(c Conn) {
	t.conn = c
}

func (t *GobWriter) Run() error {
	if t.conn == nil {
		panic("Conn not set gob writer")
	}

	encoder := gob.NewEncoder(t.conn)

	for {
		obj, ok := <-t.channel
		//fmt.Println("writer: Got Message")
		if !ok {
			//channel closed
			//fmt.Println("Writer Closed", t.conn.typ)
			t.conn.Close()
			return nil
		}

		//fmt.Println("writer: encoding Message")
		err := encoder.Encode(&obj)
		if err != nil {
			fmt.Println("Encoder error", err)
			return err
		}
	}
}

func (t *GobWriter) Channel() <-chan stream.Object {
	return t.channel
}

type GobReader struct {
	channel      chan stream.Object
	conn         Conn
	closeChannel bool
}

func NewGobReader() *GobReader {
	return &GobReader{make(chan stream.Object, 3), nil, true}
}

func (t *GobReader) SetConn(c Conn) {
	t.conn = c
}

func (t *GobReader) Run() error {
	if t.conn == nil {
		panic("Conn not set gob reader")
	}

	decoder := gob.NewDecoder(t.conn)

	if t.closeChannel {
		defer close(t.channel)
	}

	for {
		var v stream.Object
		err := decoder.Decode(&v)
		//fmt.Println("reader: Got Message")
		if err != nil {
			if err == io.EOF || t.conn.IsClosed() {
				return nil
			}
			fmt.Println("Decoder error ", t.conn.Name(), err)
			return err
		}
		//fmt.Println("reader: sent Message")
		t.channel <- v
	}
}

func (t *GobReader) SetChannel(ch chan stream.Object) {
	t.channel = ch
	t.closeChannel = false
}

func (t *GobReader) Channel() chan<- stream.Object {
	return t.channel
}
