package encoding

type ZContentType uint64

const metadataSize uint8 = 16

type Content struct {
	Type	ZContentType
	Len	uint64
	Data	[]byte
}

func ContentInit(t ZContentType, d []byte) *Content {
	return &Content {
		Type: t,
		Len: uint64(len(d)),
		Data: d,
	}
}
