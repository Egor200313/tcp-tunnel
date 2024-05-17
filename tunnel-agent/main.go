package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	tunnelIP          = flag.String("tunnel-ip", "127.0.0.1", "tunnel server addr")
	tunnelPort        = flag.Int("tunnel-port", 9999, "tunnel server port")
	localResourceIp   = flag.String("local-ip", "127.0.0.1", "local resource address")
	localResourcePort = flag.Int("local-port", 80, "local resource port")
)

const (
	BUFFER_SIZE = 100000
)

func makeAddr(ip string, port int) string {
	return ip + ":" + fmt.Sprint(port)
}

func main() {
	flag.Parse()
	tunnelAddr := makeAddr(*tunnelIP, *tunnelPort)
	tunnelConn, err := net.Dial("tcp", tunnelAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer tunnelConn.Close()

	for {
		received := make([]byte, BUFFER_SIZE)
		n, _ := tunnelConn.Read(received)
		if n > 0 {
			if string(received[:n]) == "HB" {
				tunnelConn.Write([]byte("HB"))
				continue
			}

			response := requestLocalResourse(received[:n])
			tunnelConn.Write(response)
		}
	}
}

func requestLocalResourse(request []byte) []byte {
	localResourceAddr := makeAddr(*localResourceIp, *localResourcePort)
	localConn, err := net.Dial("tcp", localResourceAddr)
	defer localConn.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = localConn.Write(request)

	if err != nil {
		log.Fatal(err)
	}

	localResponse := make([]byte, BUFFER_SIZE)
	n, err := localConn.Read(localResponse)
	if err != nil {
		log.Fatal(err)
	}

	return localResponse[:n]
}
