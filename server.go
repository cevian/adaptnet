package adaptnet

import (
	"crypto/rand"
	"fmt"

	"github.com/cevian/adaptnet/netchan"
)

type ServerOp struct {
	addr string
}

func NewServerOp(addr string) *ServerOp {
	return &ServerOp{addr}
}

func handleConnection(cp netchan.ChannelProcessor) {
	reader := cp.ChannelReader().(*netchan.ByteReader).Channel()
	writer := cp.ChannelWriter().(*netchan.ByteWriter).Channel()
	defer close(writer)

	fmt.Println("Server Getting")
	rb := <-reader
	r := &Request{}
	err := DeserializeObject(r, rb)
	if err != nil {
		panic(err)
	}

	resp := make([]byte, r.NumBytes)
	rand.Read(resp)
	fmt.Println("Server Sending")
	writer <- resp
}

func (t *ServerOp) Run() error {
	server := netchan.NewByteIndepndentServer(t.addr)
	err := server.Listen()
	if err != nil {
		panic(err)
	}

	defer server.Wait()
	defer server.Close()

	newConnCh := server.NewConnChannel()

	fmt.Println("Waiting Conn")
	cp := <-newConnCh
	fmt.Println("Launching Conn")
	handleConnection(cp)

	return nil
}

func (t *ServerOp) Stop() error {
	return nil
}
