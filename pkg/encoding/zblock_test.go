package encoding

import (
	"fmt"
	"net"
	"testing"
	"io"
)

// Test if ZBlock/Content work as expected
func TestZBlock(t *testing.T) {
	// Simulate a ZJunx server 
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
			for {
				ct := ContentInit(ZContentType(0), []byte{})
				if err := b.Unmarshalling(conn, ct); err != nil {
					if err != io.EOF {
						fmt.Println(err)
					}
					break
				}
				
				fmt.Printf("Content recved ! Type: %d; Size: %d bytes: \n%s\n", ct.Type, ct.Len, ct.Data)
			}

			ch <- true
		}
	}()

	// Simulate a ZJunx client connecting to ZJunx server
	cnx, err := net.Dial("tcp4", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	b := BlockInit()
	msg := []byte{}

	// Prepare 3 packages to send
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
