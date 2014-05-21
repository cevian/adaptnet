package adaptnet

import (
	"fmt"

	"github.com/cevian/adaptnet/netchan"
)

type ClientOp struct {
	addr string
}

func NewClientOp(addr string) *ClientOp {
	return &ClientOp{addr}
}

func (t *ClientOp) Run() error {
	client, reader, writer := netchan.NewByteClient(t.addr)
	fmt.Println("Connecting")
	err := client.Connect()
	defer client.Wait()
	defer close(writer)

	if err != nil {
		panic(err)
	}

	r := &Request{100}
	b, err := SerializeObject(r)
	fmt.Println("Sending")
	writer <- b

	response := <-reader

	if len(response) != 100 {
		panic("Wrong len")
	}

	fmt.Println("Quiting")
	return nil
}

func (t *ClientOp) Stop() error {
	return nil
}
