package server

import (
	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

type ZRequest interface {
	GetContentType() encoding.ZContentType
}

