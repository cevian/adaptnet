package adaptnet

import (
	"crypto/rand"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

/*type MemoizeMap map[int][]byte

func NewMemoizeMap() MemoizeMap {
	return make(map[int][]byte)
}

func (t *MemoizeMap) Get(size int) []byte {
	b, ok := (*t)[size]
	if !ok {
		b = make([]byte, size)
		rand.Read(b)
		(*t)[size] = b
	}
	return b
}*/

type PayloadGen struct {
	basis []byte
}

func NewPayloadGen(max_size int) *PayloadGen {
	b := make([]byte, max_size)
	rand.Read(b)
	return &PayloadGen{b}
}

func (t *PayloadGen) Write(w io.Writer, size int) error {
	count := 0
	for count < size {
		left := size - count
		to_send := left
		if to_send > len(t.basis) {
			to_send = len(t.basis)
		}
		_, err := w.Write(t.basis[:to_send])
		if err != nil {
			panic(err)
			return err
		}
		count += to_send
	}
	return nil
}

func (t *PayloadGen) Get(size int) []byte {
	if size > len(t.basis) {
		panic("size above maxsize")
	}
	b := t.basis[:size]
	return b
}

type ServerOp struct {
	addr           string
	numConnections int
	maxPayload     int
}

func NewServerOp(addr string, numConnections int, maxPayload int) *ServerOp {
	return &ServerOp{addr, numConnections, maxPayload}
}

func handleConnection(cp netchan.ChannelProcessor, syncCh chan bool, connNo int, pg *PayloadGen) {
	reader := cp.ChannelReader().(*netchan.ByteReader).Channel()
	writer := cp.ChannelWriter().(*netchan.ByteWriter).Channel()
	defer close(writer)

	<-syncCh
	fmt.Println("Server Started", connNo)
	for rb := range reader {
		r := &Request{}
		err := DeserializeObject(r, rb)
		if err != nil {
			panic(err)
		}

		start := time.Now()
		resp := pg.Get(int(r.NumBytes))
		fmt.Println("Server Sending", connNo, time.Since(start))
		writer <- resp
	}
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
	syncCh := make(chan bool, 0)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	pg := NewPayloadGen(t.maxPayload)
	for i := 0; i < t.numConnections; i++ {
		cp := <-newConnCh
		connNo := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(cp, syncCh, connNo, pg)
		}()
	}
	close(syncCh)
	fmt.Println("Done getting connections")

	return nil
}

func (t *ServerOp) Stop() error {
	return nil
}
