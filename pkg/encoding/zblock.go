package encoding

import (
	"bytes"
	"encoding/binary"
)

type ZBlock interface {
	Marshalling(ZContent)([]byte, error)
	Unmarshalling([]byte, ZContent) error
}

type Block struct {}

func BlockInit() ZBlock {
	return &Block{}
}

func (b *Block) Marshalling(ct ZContent)([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	data := []interface{}{
		ct.(Content).Type,
		ct.(Content).Len,
		ct.(Content).Data,
	}

	for _,v := range data {
		if err := binary.Write(buf, binary.BigEndian, v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (b *Block) Unmarshalling(data []byte, ct ZContent) error {
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, ct.(Content).Type); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, ct.(Content).Len); err != nil {
		return err
	}

	return nil
}
