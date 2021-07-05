package main

import (
	"net"
	"log"
	"time"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	blk := encoding.BlockInit()
	
	for i:=7; i<10; i++ {
		req, err := blk.Marshalling(encoding.ContentInit(encoding.ZContentType(i), []byte("Hello ZJunx!\n")))
		if err != nil {
			log.Fatalln("Marshalling error: ", err)
		}

		if _, err := conn.Write(req); err != nil {
			log.Fatalln("Write to conn error: ", err)
		}
		log.Println("Request sent.")

		resp := encoding.ContentInit(encoding.ZContentType(0), []byte{})
		if err := blk.Unmarshalling(conn, resp); err != nil {
			log.Fatalln("Unmarshalling error: ", err)
		}

		log.Printf("Get response from server: [%d] %s\n", resp.Type, resp.Data)
		
		time.Sleep( 5 * time.Second )
	}
}
