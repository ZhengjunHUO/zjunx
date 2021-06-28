package main

import (
	"net"
	"log"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	
	for i:=0; i<3; i++ {
		if _, err := conn.Write([]byte("Hello ZJunx!\n")); err != nil {
			log.Println(err)
		}
		
		time.Sleep( 10 * time.Second )	
	}
}
