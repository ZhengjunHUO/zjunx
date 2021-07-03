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

	blk, ct := encoding.BlockInit(), encoding.ContentInit(encoding.ZContentType(8), []byte("Hello ZJunx!\n"))
	
	for i:=0; i<3; i++ {
		req, err := blk.Marshalling(ct)
		if err != nil {
			log.Println(err)
		}

		if _, err := conn.Write(req); err != nil {
			log.Println(err)
		}
		log.Println("Request sent.")

		resp := encoding.ContentInit(encoding.ZContentType(0), []byte{})
		if err := blk.Unmarshalling(conn, resp); err != nil {
			log.Println(err)
		}

		log.Printf("Get response from server: [%d]%s\n", resp.Type, resp.Data)
		
		time.Sleep( 10 * time.Second )
	}
}
