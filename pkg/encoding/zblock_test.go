package encoding

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestZBlock(t *testing.T) {
	l, err := net.Listen("tcp4", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan bool)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}

			b := BlockInit()
			metaData := make([]byte, metadataSize)
			for {
				if _, err := io.ReadFull(conn, metaData); err != nil {
					fmt.Println(err)
					break
				}

				ct := ContentInit(ZContentType(0), []byte{})
				if err := b.Unmarshalling(metaData, ct); err != nil {
					fmt.Println(err)
					break
				}

				if ct.Len > 0 {
					ct.Data = make([]byte, ct.Len)
					if _, err := io.ReadFull(conn, ct.Data); err != nil {
						fmt.Println(err)
						return
					}
				
					fmt.Printf("Content recved ! Type: %d; Size: %d bytes: \n%s\n", ct.Type, ct.Len, ct.Data)
				}
			}

			ch <- true
		}
	}()

	cnx, err := net.Dial("tcp4", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	b := BlockInit()
	msg := []byte{}

	ct := make([]*Content, 3)
	ct[0] = ContentInit(ZContentType(3), []byte("Thanks for"))
	ct[1] = ContentInit(ZContentType(2), []byte("using "))
	ct[2] = ContentInit(ZContentType(5), []byte("Zjunx !"))

	for i := range ct {
		d, err := b.Marshalling(ct[i])
		if err != nil {
			fmt.Println(err)
			return
		}

		msg = append(msg, d...)
	}

	if _, err := cnx.Write(msg); err != nil {
		fmt.Println(err)
		return
	}	

	cnx.Close()

	select{
	case <-ch:
		fmt.Println("Connection closed from client side ...")
	}
}
