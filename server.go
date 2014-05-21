package adaptnet

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/cevian/adaptnet/netchan"
)

type MemoizeMap map[int][]byte

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
}

type ServerOp struct {
	addr           string
	numConnections int
}

func NewServerOp(addr string, numConnections int) *ServerOp {
	return &ServerOp{addr, numConnections}
}

func handleConnection(cp netchan.ChannelProcessor, syncCh chan bool, connNo int) {
	reader := cp.ChannelReader().(*netchan.ByteReader).Channel()
	writer := cp.ChannelWriter().(*netchan.ByteWriter).Channel()
	defer close(writer)

	mm := NewMemoizeMap()
	<-syncCh
	fmt.Println("Server Started", connNo)
	for rb := range reader {
		r := &Request{}
		err := DeserializeObject(r, rb)
		if err != nil {
			panic(err)
		}

		start := time.Now()
		resp := mm.Get(int(r.NumBytes))
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

	for i := 0; i < t.numConnections; i++ {
		cp := <-newConnCh
		connNo := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(cp, syncCh, connNo)
		}()
	}
	close(syncCh)
	fmt.Println("Done getting connections")

	return nil
}

func (t *ServerOp) Stop() error {
	return nil
}
