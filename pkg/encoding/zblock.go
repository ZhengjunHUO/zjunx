package encoding

type ZBlock interface {
	Marshalling(ZContent)([]byte, error)
	Unmarshalling([]byte, ZContent) error
}
