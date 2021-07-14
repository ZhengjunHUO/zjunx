package main

import (
	"net"
	"log"
	"bufio"
	"os"
	"fmt"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

func main() {
	var username string = "福福"

	conn, err := net.Dial("tcp4", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	blk := encoding.BlockInit()
	
	req, err := blk.Marshalling(encoding.ContentInit(encoding.ZContentType(1), []byte(username)))
	if err != nil {
		log.Fatalln("Marshalling error: ", err)
	}

	if _, err := conn.Write(req); err != nil {
		log.Fatalln("Write to conn error: ", err)
	}

	resp := encoding.ContentInit(encoding.ZContentType(0), []byte{})
	if err := blk.Unmarshalling(conn, resp); err != nil {
		log.Fatalln("Unmarshalling error: ", err)
	}

	log.Printf("%s\n", resp.Data)

	go func() {
		for {
			resp := encoding.ContentInit(encoding.ZContentType(0), []byte{})
			if err := blk.Unmarshalling(conn, resp); err != nil {
				log.Fatalln("Unmarshalling error: ", err)
			}

			fmt.Printf("%s\n", resp.Data)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		req, err = blk.Marshalling(encoding.ContentInit(encoding.ZContentType(2), scanner.Bytes() ))
		if err != nil {
			log.Fatalln("Marshalling error: ", err)
		}

		if _, err := conn.Write(req); err != nil {
			log.Fatalln("Write to conn error: ", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
