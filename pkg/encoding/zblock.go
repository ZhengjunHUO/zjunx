package encoding

import (
	"bytes"
	"encoding/binary"
)

type ZBlock interface {
	Marshalling(*Content)([]byte, error)
	Unmarshalling([]byte, *Content) error
}

type Block struct {}

func BlockInit() ZBlock {
	return &Block{}
}

func (b *Block) Marshalling(ct *Content)([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	data := []interface{}{
		ct.Type,
		ct.Len,
		ct.Data,
	}

	for _,v := range data {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (b *Block) Unmarshalling(data []byte, ct *Content) error {
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, &ct.Type); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &ct.Len); err != nil {
		return err
	}

	return nil
}
