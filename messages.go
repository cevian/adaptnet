package adaptnet

import (
	"encoding/binary"
	"io"
)

type Request struct {
	NumBytes int32
}

func (t *Request) Serialize(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, &t.NumBytes); err != nil {
		return err
	}
	return nil
}

func (t *Request) Deserialize(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &t.NumBytes); err != nil {
		return err
	}
	return nil
}

type Response struct {
	Payload []byte
}
