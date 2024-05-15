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
	conn, err := net.Dial("tcp", strings.Join([]string{*remoteAddr, *remotePort}, ":"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localConn, err := net.Dial("tcp", *localResource)
	if err != nil {
		log.Fatal(err)
	}
	defer localConn.Close()

	for {
		received := make([]byte, 1024)
		n, err := conn.Read(received)
		//received, err := io.ReadAll(conn)
		if err != nil {
			//log.Println(err)
		}
		if n > 0 {
			//log.Println(received, "from server")
			_, err := localConn.Write(received[:n])
			if err != nil {
				log.Fatal(err)
			}

			localResponse := make([]byte, 1024)
			n, err := localConn.Read(localResponse)
			if err != nil {
				log.Fatal(err)
			}
			//localResponse, err := io.ReadAll(localConn)
			//log.Println("from local server response", localResponse[:n])
			conn.Write(localResponse[:n])
		}
	}
}
