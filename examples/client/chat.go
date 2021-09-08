package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"bufio"
	"os"
	"fmt"

	"github.com/ZhengjunHUO/zjunx/pkg/encoding"
)

const caCrt = `-----BEGIN CERTIFICATE-----
MIIC7zCCAdegAwIBAgIJAKo1cTNaEVy2MA0GCSqGSIb3DQEBDQUAMA4xDDAKBgNV
BAMMA2h1bzAeFw0yMTA5MDgwNzMzMDBaFw0zMTA5MDYwNzMzMDBaMA4xDDAKBgNV
BAMMA2h1bzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALNHE7GHOL/b
J9+Lr23Ew7bWhYGoUzEYG/pccNUilq5GCQzn0/H5U4Qy6bZs9qXDSRkAg1DWRwJD
7k4Pczmo9zTz2U5wUNeJlHrJruE57ehaZMUvy6U/5e/8FG8FKir3VSUxcNyO4FzJ
MiiN604DMCdTXO+am+tNkjxx99iUuN92OwVOni0RIPIcPb9Nn4Brt2wO5aEZ8ucl
L5PKx4quSi05o32W0WvCjz+70ytrvyNtRutyEjXNVs9Nd+lypdZ/oSBqzvgJ6JVQ
oqwhHDAzjD53d3hQyRFBWUAVE7Y7M55oSHtWfwfZdMR+i6JRDC4eC7cA3YFbHCdS
fCjYQQnzkOkCAwEAAaNQME4wHQYDVR0OBBYEFHvnasQDiXtGoPI7TTTjKR7sRACW
MB8GA1UdIwQYMBaAFHvnasQDiXtGoPI7TTTjKR7sRACWMAwGA1UdEwQFMAMBAf8w
DQYJKoZIhvcNAQENBQADggEBAD6SAD49fZPUnF6Updp6MICGs5+RAp/WIlo3IMMP
O0ChiAwjq/hhZXGsO3iA/LmUyGi7ZeQca18zQFD5BZjqCTj1YSJv/wIsHBE7zDAE
1JIkoV2Cv8b7Sx4Y2S7y+8xEUpC9BzDfwjNlCGzQ1OztHWKA8ImxewvJYqdFmqL9
IlZ+e4mxEDmNMBxKnVOW1xRQrggfBupLI2DszegGRUaoNYxMM8yW4ZkwtS091wF2
FfIZJrI15Icx5jEgl2l/KC6AmmVd3RFqTBXlIUaCZgURB+9nxVxdpTE5+xyKYWi6
74WnwuJL3ybG44EDlQqjsXRShDjHJQcMANMdVSi3c1qbcro=
-----END CERTIFICATE-----`

func main() {
	var username string = "福福"

	certs := x509.NewCertPool()
	if ok := certs.AppendCertsFromPEM([]byte(caCrt)); !ok {
		log.Fatal("Failed to load ca crt !")
	}
	config := &tls.Config{RootCAs: certs, ServerName: "localhost"}

	conn, err := tls.Dial("tcp4", "127.0.0.1:8080", config)
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
