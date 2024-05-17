package main

import (
	"flag"
	"log"
	"net"
	"strings"
)

var (
	remoteAddr    = flag.String("remote-addr", "localhost", "tunnel server addr")
	remotePort    = flag.String("remote-port", "9999", "tunnel server port")
	localResource = flag.String("addr", "localhost:80", "local resource address")
)

func main() {
	flag.Parse()
	tunnelConn, err := net.Dial("tcp", strings.Join([]string{*remoteAddr, *remotePort}, ":"))
	if err != nil {
		log.Fatal(err)
	}
	defer tunnelConn.Close()

	for {
		received := make([]byte, 1024)
		n, _ := tunnelConn.Read(received)
		if n > 0 {
			if string(received[:n]) == "HB" {
				//log.Println("HB received")
				//m, err :=
				tunnelConn.Write([]byte("HB"))
				//log.Println("sending to tunnel HB:", m, err)
				continue
			}

			response := requestLocalResourse(received[:n])
			tunnelConn.Write(response)
		}
	}
}

func requestLocalResourse(request []byte) []byte {
	localConn, err := net.Dial("tcp", *localResource)
	defer localConn.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = localConn.Write(request)

	if err != nil {
		log.Fatal(err)
	}

	localResponse := make([]byte, 1024)
	n, err := localConn.Read(localResponse)
	if err != nil {
		log.Fatal(err)
	}

	return localResponse[:n]
}
