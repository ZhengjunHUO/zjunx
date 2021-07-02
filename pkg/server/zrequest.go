package server

import (
	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

type ZRequest interface {
	GetContentType() encoding.ZContentType
}

type Request struct {
	conn	ZConnection
	cont	*encoding.Content
}

func ReqInit(conn ZConnection, cont *encoding.Content) ZRequest {
	return &Request{
		conn: conn,
		cont: cont,
	}
}

func (r *Request) GetContentType() encoding.ZContentType {
	return r.cont.Type
}
