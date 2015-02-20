package adaptnet

import (
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

type ServerDirectOp struct {
	addr           string
	numConnections int
	maxPayload     int
}

func NewServerDirectOp(addr string, numConnections int, maxPayload int) *ServerDirectOp {
	return &ServerDirectOp{addr, numConnections, maxPayload}
}

func handleDirectConnection(cp netchan.ChannelProcessor, syncCh chan bool, connNo int, pg *PayloadGen) {
	reader := cp.ChannelReader().(*netchan.ByteReader)
	writer := cp.ChannelWriter().(*netchan.ByteWriter)
	defer writer.CloseConn()

	<-syncCh
	fmt.Println("Server Started", connNo)
	for {
		before_request := time.Now()
		rb, err := reader.ReadConnection()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		request_pause := time.Since(before_request)

		r := &Request{}
		err = DeserializeObject(r, rb)
		if err != nil {
			panic(err)
		}

		start := time.Now()
		//resp := pg.Get(int(r.NumBytes))
		//fmt.Println("Server Sending On Connenction", connNo, " Time to generate: ", time.Since(start), "Time before request", request_pause)
		//err = writer.WriteConnection(resp)
		w, err := writer.GetWriter(int(r.NumBytes))
		if err != nil {
			panic(err)
		}
		pg.Write(w, int(r.NumBytes))
		fmt.Println("Server Sending On Connenction", connNo, " Time to generate and send: ", time.Since(start), "Time before request", request_pause)

		if err != nil {
			panic(err)
		}
	}
}

func (t *ServerDirectOp) Run() error {
	server := netchan.NewByteIndepndentServer(t.addr)
	server.AutoStart = false
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
		fmt.Println("Waiting for connection");
		cp := <-newConnCh
		connNo := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleDirectConnection(cp, syncCh, connNo, pg)
		}()
	}
	close(syncCh)
	fmt.Println("Done getting connections")

	return nil
}

func (t *ServerDirectOp) Stop() error {
	return nil
}
