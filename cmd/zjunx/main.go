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

	msg := encoding.ContentInit(encoding.ZContentType(2), []byte(username+"上线了🎉"))
	req.Connection().GetServer().GetMux().Handle(server.ReqInit(req.Connection(), msg))
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

func AnnonceOffline(conn server.ZConnection) {
	if username := conn.GetContext("Username"); username != nil {
		msg := encoding.ContentInit(encoding.ZContentType(2), []byte(string(username.(string))+"下线了👋"))
		conn.GetServer().GetMux().Handle(server.ReqInit(conn, msg))
	}
}

func main() {
	server.NewServer(
		server.WithHandler(encoding.ZContentType(1), &LoginHandler{}),
		server.WithHandler(encoding.ZContentType(2), &BroadcastHandler{}),
		server.WithHandler(encoding.ZContentType(8), &EchoHandler{}),
		server.WithPreStop(AnnonceOffline),
	).Start()
}
