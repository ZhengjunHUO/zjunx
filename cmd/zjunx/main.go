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

type LoginHandler struct {
	server.ZHandler
}

func (lh *LoginHandler) Handle(req server.ZRequest) {
	cid := req.Connection().GetID()
	username := string(req.ContentData())
	log.Printf("[DEBUG] [Connection %v] LoginHandler receive data: %s\n", cid, username)
	req.Connection().UpdateContext("Username", username)
	log.Printf("[DEBUG] [Connection %v] context upgrade.\n", cid)
	str := "你好"+username+"! 欢迎使用ZJunx!"
	if err := req.Connection().RespondToClient(encoding.ZContentType(24), []byte(str)); err != nil {
		log.Println("[DEBUG] EchoHandler error: ", err)
	}
}

type BroadcastHandler struct {
	server.ZHandler
}

func (bh *BroadcastHandler) Handle(req server.ZRequest) {
	cid := req.Connection().GetID()
	log.Printf("[DEBUG] [Connection %v] BroadcastHandler receive data: %s\n", cid, string(req.ContentData()))

	var resp []byte
	if val := req.Connection().GetContext("Username"); val != nil {
		v := string(val.(string))
		v += ": "
		resp = append([]byte(v), req.ContentData()...)
	}else{
		resp = req.ContentData()
	}

	for _, cnx := range req.Connection().GetServer().GetCnxAdm().RetrieveAll() {
		cnx.RespondToClient(encoding.ZContentType(24), resp)
		log.Printf("[DEBUG] [Connection %v] Broadcast to connection %v\n", cid, cnx.GetID())
	}
}

func main() {
	s := server.ServerInit()
	s.RegistHandler(encoding.ZContentType(1), &LoginHandler{})
	s.RegistHandler(encoding.ZContentType(2), &BroadcastHandler{})
	s.RegistHandler(encoding.ZContentType(8), &EchoHandler{})
	s.Start()
}
