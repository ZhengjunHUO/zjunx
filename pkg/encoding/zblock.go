package encoding

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ZBlock interface {
	Marshalling(*Content)([]byte, error)
	Unmarshalling(io.Reader, *Content) error
}

type Block struct {}

func BlockInit() ZBlock {
	return &Block{}
}

// Serialize the Content struct to raw bytes
func (b *Block) Marshalling(ct *Content)([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	data := []interface{}{
		ct.Type,
		ct.Len,
		ct.Data,
	}

	// Write to a buffer in order (Type, Len, Data)
	for _,v := range data {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Decode the raw bytes (TCP stream) to Content struct
func (b *Block) Unmarshalling(conn io.Reader, ct *Content) error {
	// Read the header part (Type, Len) of the TCP payload
	meta := make([]byte, metadataSize)
	if _, err := io.ReadFull(conn, meta); err != nil {
		return err
	}

	r := bytes.NewReader(meta)
	if err := binary.Read(r, binary.BigEndian, &ct.Type); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &ct.Len); err != nil {
		return err
	}

	// Read the data part (according to its Len given in metadata) of the TCP payload
	if ct.Len > 0 {
		ct.Data = make([]byte, ct.Len)
		if _, err := io.ReadFull(conn, ct.Data); err != nil {
			return err
		}
	}

	return nil
}
