package main

import (
	"flag"
	"log"
	"net"
)

var (
	remoteIP          = flag.String("tunnel-ip", "127.0.0.1", "tunnel server addr")
	remotePort        = flag.String("tunnel-port", "9999", "tunnel server port")
	localResourceIp   = flag.String("local-ip", "127.0.0.1", "local resource address")
	localResourcePort = flag.String("local-port", "80", "local resource port")
)

func main() {
	flag.Parse()
	tunnelAddr := *remoteIP + ":" + *remotePort
	tunnelConn, err := net.Dial("tcp", tunnelAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer tunnelConn.Close()

	for {
		received := make([]byte, 1024)
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
	localResourceAddr := *localResourceIp + ":" + *localResourcePort
	localConn, err := net.Dial("tcp", localResourceAddr)
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
