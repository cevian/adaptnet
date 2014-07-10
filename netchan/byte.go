package netchan

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"time"
)

type ByteWriterCollection struct {
	last_idx int
	coll     map[int]*ByteWriter
}

func NewByteWriterCollection() *ByteWriterCollection {
	return &ByteWriterCollection{0, make(map[int]*ByteWriter)}
}

func (t *ByteWriterCollection) Add(bw *ByteWriter) int {
	t.last_idx++
	t.coll[t.last_idx] = bw
	return t.last_idx
}

func (t *ByteWriterCollection) ByteWriter(idx int) *ByteWriter {
	return t.coll[idx]
}

func (t *ByteWriterCollection) Remove(idx int) {
	delete(t.coll, idx)
}

type ByteDemuxProcessorGenerator struct {
	DemuxChannel chan *DemuxByteWrapper
	bcc          *ByteWriterCollection
}

func (t *ByteDemuxProcessorGenerator) NewChannelProcessor() ChannelProcessor {
	writer := NewByteWriter()
	nid := t.bcc.Add(writer)
	br := NewDemuxByteReader(nid, t.DemuxChannel)
	closefn := func(id int) {
		close(t.bcc.ByteWriter(id).Channel())
		t.bcc.Remove(id)
	}
	br.CloseCallback = closefn
	return &SimpleChannelProcessor{&sync.WaitGroup{}, br, writer}
}

type ByteIndependentProcessorGenerator struct {
}

func (t *ByteIndependentProcessorGenerator) NewChannelProcessor() ChannelProcessor {
	br := NewByteReader()
	bw := NewByteWriter()
	p := &SimpleChannelProcessor{&sync.WaitGroup{}, br, bw}
	return p
}

func NewByteClient(addr string) (*Client, <-chan []byte, chan<- []byte) {
	br := NewByteReader()
	bw := NewByteWriter()
	c := NewClient(addr, &SimpleChannelProcessor{&sync.WaitGroup{}, br, bw})
	return c, br.Channel(), bw.Channel()
}

func NewByteDemuxServer(addr string) (*Server, <-chan *DemuxByteWrapper, *ByteWriterCollection) {
	ch := make(chan *DemuxByteWrapper, 0)
	bcc := NewByteWriterCollection()
	bpg := &ByteDemuxProcessorGenerator{ch, bcc}
	return NewServer(addr, bpg), ch, bcc
}

func NewByteIndepndentServer(addr string) *Server {
	bpg := &ByteIndependentProcessorGenerator{}
	return NewServer(addr, bpg)
}

type ByteWriter struct {
	channel chan []byte
	conn    Conn
}

func NewByteWriter() *ByteWriter {
	return &ByteWriter{make(chan []byte, 0), nil}
}

func (t *ByteWriter) SetConn(c Conn) {
	t.conn = c
}

func (t *ByteWriter) CloseConn() {
	t.conn.Close()
}

func (t *ByteWriter) Run() error {
	if t.conn == nil {
		panic("Conn not set byte writer")
	}

	for {
		obj, ok := <-t.channel
		//fmt.Println("writer: Got Message")
		if !ok {
			//channel closed
			//fmt.Println("Writer Closed", t.conn.typ)
			t.conn.Close()
			return nil
		}

		err := t.WriteConnection(obj)
		if err != nil {
			return err
		}
	}
}

func (t *ByteWriter) GetWriter(length_will_send int) (io.Writer, error) {
	err := binary.Write(t.conn, binary.LittleEndian, uint32(length_will_send))
	if err != nil {
		fmt.Println("write length", err)
		return nil, err
	}
	return t.conn, nil
}

func (t *ByteWriter) WriteConnection(obj []byte) error {
	w, err := t.GetWriter(len(obj))
	if err != nil {
		return err
	}
	_, err = w.Write(obj)
	if err != nil {
		fmt.Println("write error", err)
		return err
	}
	return nil
}

func (t *ByteWriter) Channel() chan<- []byte {
	return t.channel
}

type ByteReader struct {
	channel      chan []byte
	conn         Conn
	closeChannel bool
}

func NewByteReader() *ByteReader {
	return &ByteReader{make(chan []byte, 0), nil, true}
}

func (t *ByteReader) SetConn(c Conn) {
	t.conn = c
}

func (t *ByteReader) readError(err error) error {
	if err == io.EOF || t.conn.IsClosed() {
		return nil
	}
	fmt.Println("Decoder error ", t.conn.Name(), err)
	return err
}

func (t *ByteReader) Run() error {
	if t.conn == nil {
		panic("Conn not set gob reader")
	}

	if t.closeChannel {
		defer close(t.channel)
	}

	for {

		b, err := t.ReadConnection()
		if err != nil {
			return t.readError(err)
		}
		t.channel <- b
	}
}

func (t *ByteReader) ReadConnection() ([]byte, error) {
	var length uint32
	err := binary.Read(t.conn, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	b := make([]byte, length)
	_, err = io.ReadFull(t.conn, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (t *ByteReader) ReadConnectionInto(b []byte) ([]byte, time.Time, error) {
	var length uint32
	err := binary.Read(t.conn, binary.LittleEndian, &length)
	start := time.Now()
	if err != nil {
		return nil, start, err
	}

	if cap(b) < int(length) {
		b = make([]byte, length)
	}
	b = b[:length]

	_, err = io.ReadFull(t.conn, b)
	if err != nil {
		return nil, start, err
	}
	return b, start, nil
}

func (t *ByteReader) ReadConnectionIntoNoExpand(b []byte) (time.Time, int, error) {
	var lengthU uint32
	err := binary.Read(t.conn, binary.LittleEndian, &lengthU)
	start := time.Now()
	if err != nil {
		return start, 0, err
	}

	length := int(lengthU)
	lengthRead := 0
	for lengthRead < length {
		toRead := cap(b)
		left := length - lengthRead
		if left < toRead {
			toRead = left
		}

		bufRound := b[:toRead]
		_, err = io.ReadFull(t.conn, bufRound)
		if err != nil {
			return start, 0, err
		}
		lengthRead += len(b)
	}

	return start, length, nil
}

func (t *ByteReader) SetChannel(ch chan []byte) {
	t.channel = ch
	t.closeChannel = false
}

func (t *ByteReader) Channel() <-chan []byte {
	return t.channel
}

type DemuxByteWrapper struct {
	id   int
	data []byte
}

type DemuxByteReader struct {
	channel       chan *DemuxByteWrapper
	conn          Conn
	CloseCallback func(id int)
	id            int
}

func NewDemuxByteReader(id int, ch chan *DemuxByteWrapper) *DemuxByteReader {
	return &DemuxByteReader{ch, nil, nil, id}
}

func (t *DemuxByteReader) SetConn(c Conn) {
	t.conn = c
}

func (t *DemuxByteReader) readError(err error) error {
	if err == io.EOF || t.conn.IsClosed() {
		return nil
	}
	fmt.Println("Decoder error ", t.conn.Name(), err)
	return err
}

func (t *DemuxByteReader) Run() error {
	if t.conn == nil {
		panic("Conn not set gob reader")
	}

	if t.CloseCallback != nil {
		defer t.CloseCallback(t.id)
	}

	for {
		var length uint32
		err := binary.Read(t.conn, binary.LittleEndian, &length)
		if err != nil {
			return t.readError(err)
		}

		b := make([]byte, length)
		_, err = io.ReadFull(t.conn, b)
		if err != nil {
			return t.readError(err)
		}
		t.channel <- &DemuxByteWrapper{t.id, b}
	}
}

func (t *DemuxByteReader) Channel() <-chan *DemuxByteWrapper {
	return t.channel
}
