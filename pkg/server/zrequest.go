package server

import (
	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

type ZRequest interface {
	ContentType() encoding.ZContentType
	ContentData() []byte
	Connection()  ZConnection
}

type Request struct {
	Conn	ZConnection
	Cont	*encoding.Content
}

func ReqInit(conn ZConnection, cont *encoding.Content) ZRequest {
	return &Request{
		Conn: conn,
		Cont: cont,
	}
}

func (r *Request) ContentType() encoding.ZContentType {
	return r.Cont.Type
}

func (r *Request) ContentData() []byte {
	return r.Cont.Data
}

func (r *Request) Connection() ZConnection {
	return r.Conn
}
