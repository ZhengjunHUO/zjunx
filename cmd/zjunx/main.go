package main

import (
	"github.com/ZhengjunHUO/zjunx/pkg/server"
)

func main() {
	s := server.ServerInit()
	s.Start()
}
