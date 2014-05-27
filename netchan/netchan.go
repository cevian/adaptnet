package netchan

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type ChannelReader interface {
	SetConn(c Conn)
	Run() error
}

type ChannelWriter interface {
	SetConn(c Conn)
	Run() error
}

//The way to close a channel processor is to close the ChannelWriter.
//That should close everything started by Start() and free up the WaitGroup
//A eof closes the ChannelReader channel. This should be supported better where
//the server may share on ChannelReader by multiple clients
type ChannelProcessor interface {
	Start() error
	SetWaitGroup(wg *sync.WaitGroup)
	Wait()
	ChannelReader() ChannelReader
	ChannelWriter() ChannelWriter
	SetConn(c Conn)
}

type ChannelProcessorGenerator interface {
	NewChannelProcessor() ChannelProcessor
}

type SimpleChannelProcessor struct {
	wg *sync.WaitGroup
	cr ChannelReader
	cw ChannelWriter
}

func (t *SimpleChannelProcessor) SetWaitGroup(wg *sync.WaitGroup) {
	t.wg = wg
}

func (t *SimpleChannelProcessor) Start() error {
	t.wg.Add(2)
	go func() {
		defer t.wg.Done()
		if err := t.cw.Run(); err != nil {
			fmt.Println("Writer quit err = ", err)
		}
	}()
	go func() {
		defer t.wg.Done()
		if err := t.cr.Run(); err != nil {
			fmt.Println("Reader quit err = ", err)
		}
	}()
	return nil
}

func (t *SimpleChannelProcessor) Wait() {
	t.wg.Wait()
}

func (t *SimpleChannelProcessor) ChannelWriter() ChannelWriter {
	return t.cw
}

func (t *SimpleChannelProcessor) ChannelReader() ChannelReader {
	return t.cr
}

func (t *SimpleChannelProcessor) SetConn(conn Conn) {
	t.cw.SetConn(conn)
	t.cr.SetConn(conn)
}

//convenience interface to be used by ChannelReader and ChannelWriter
type Conn interface {
	io.Reader
	io.Writer
	Close() error
	IsClosed() bool
	Name() string
}

type NetConn struct {
	net.Conn
	isClosed bool
	name     string
}

func NewNetConn(c net.Conn, name string) Conn {
	return &NetConn{c, false, name}
}

func (t *NetConn) Close() error {
	t.isClosed = true
	t.Conn.Close()
	return nil
}

func (t *NetConn) IsClosed() bool {
	return t.isClosed
}

func (t *NetConn) Name() string {
	return t.Name()
}

type Client struct {
	addr      string
	processor ChannelProcessor
	AutoStart bool
}

func NewClient(addr string, processor ChannelProcessor) *Client {
	return &Client{addr, processor, true}
}

func (t *Client) Connect() error {
	conn, err := net.Dial("tcp", t.addr)
	if err != nil {
		fmt.Println("Cannot establish a connection with %s %v", t.addr, err)
		return err
	}

	c := NewNetConn(conn, "Client")
	t.processor.SetConn(c)
	if t.AutoStart {
		t.processor.Start()
	}

	return nil
}

func (t *Client) Processor() ChannelProcessor {
	return t.processor
}

func (t *Client) Wait() {
	t.Processor().Wait()
}

type Server struct {
	addr           string
	gen            ChannelProcessorGenerator
	newConnChannel chan ChannelProcessor
	listener       net.Listener
	closing        bool
	wg             *sync.WaitGroup
	AutoStart      bool
}

func NewServer(addr string, gen ChannelProcessorGenerator) *Server {
	return &Server{addr, gen, make(chan ChannelProcessor, 3), nil, false, &sync.WaitGroup{}, true}
}

func (t *Server) NewConnChannel() chan ChannelProcessor {
	return t.newConnChannel
}

func (t *Server) Listen() error {
	ln, err := net.Listen("tcp", t.addr)
	if err != nil {
		fmt.Println("Error listening %v", err)
		return err
	}
	t.listener = ln

	t.wg.Add(1)
	f := func() {
		defer t.wg.Done()
		for {
			conn, err := t.listener.Accept()
			if err != nil {
				if !t.closing {
					fmt.Println("Error Accepting %v", err)
					return
				}
				return
			}

			cp := t.gen.NewChannelProcessor()
			c := NewNetConn(conn, "Server")
			cp.SetConn(c)
			cp.SetWaitGroup(t.wg)
			if t.AutoStart {
				cp.Start()
			}
			t.newConnChannel <- cp
		}
	}

	go f()
	return nil
}

func (t *Server) Close() error {
	t.closing = true
	return t.listener.Close()
}

func (t *Server) Wait() {
	t.wg.Wait()
}

/*
func (t *Writer) Encode(obj interface{}) ([]byte, error) {
	return disttopk.GobBytesEncode(obj)
}

func (t Writer) writeMsg(msg []byte) error {
	err := binary.Write(t.writer, binary.LittleEndian, uint32(len(msg)))
	if err != nil {
		return err
	}
	_, err = t.writer.Write(msg)
	return err
}*/
