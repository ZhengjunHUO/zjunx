package encoding

type ZContentType uint64

const metadataSize uint8 = 16

type ZContent interface {}

type Content struct {
	Type	ZContentType
	Len	uint64
	Data	[]byte
}

func ContentInit(t ZContentType, d []byte) ZContent {
	return &Content {
		Type: t,
		Len: uint64(len(d)),
		Data: d,
	}
}
