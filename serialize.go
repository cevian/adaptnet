package adaptnet

import (
	"bytes"
	"io"
)

type Serializer interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
}

func SerializeObject(obj Serializer) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := obj.Serialize(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeserializeObject(into Serializer, b []byte) error {
	buf := bytes.NewReader(b)
	return into.Deserialize(buf)
}
