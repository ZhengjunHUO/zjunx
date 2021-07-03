package main

import (
	"log"
	"github.com/ZhengjunHUO/zjunx/pkg/server"
	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

type EchoHandler struct {
	server.ZHandler
}

func (ec *EchoHandler) Handle(req server.ZRequest) {
	log.Println("[DEBUG] EchoHandler receive data: ", string(req.ContentData()))
	if err := req.Connection().RespondToClient(encoding.ZContentType(24), req.ContentData()); err != nil {
		log.Println("[DEBUG] EchoHandler error: ", err)
	}
}

func main() {
	s := server.ServerInit()
	s.RegistHandler(encoding.ZContentType(8), &EchoHandler{})
	s.Start()
}
