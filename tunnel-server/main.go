package main

import (
	"flag"
	"log"
	"net"
	"time"
)

var (
	clientPort = flag.Int("client-port", 9090, "local server port for public clients requests")
	tunnelPort = flag.Int("tunnel-port", 9999, "local server port to accept tunnel agent connections")
	tunnelConn = &net.TCPConn{}
	clientConn = &net.TCPConn{}

	clientIsResponded = make(chan struct{}, 1)
)

func main() {
	flag.Parse()
	setTunnelConn()
	clientLoop()
}

func heartBeatsRoutine(conn *net.TCPConn) {
	for {
		_, err := conn.Write([]byte("HB"))
		if err != nil {
			return
		}
		time.Sleep(3 * time.Second)
	}
}

func clientLoop() {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *clientPort})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("waiting for client...")
	for {
		clientConn, err = l.AcceptTCP()
		//conn.SetKeepAlive(true)
		//conn.SetKeepAlivePeriod(time.Second * 30)

		if err != nil {
			log.Fatalf("accepting client connection: %s\n", err)
		}
		log.Println("got client connection")
		transferClientRequestToTunnel()
	}
}

func setTunnelConn() {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *tunnelPort})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("waiting for agent...")
	tunnelConn, err = l.AcceptTCP()
	if err != nil {
		log.Fatalf("accepting agent connection: %s\n", err)
	}
	log.Println("got agent connection")
	go heartBeatsRoutine(tunnelConn)
	go receiver()
}

func transferClientRequestToTunnel() {
	defer clientConn.Close()

	receivedFromClient := make([]byte, 1024)
	n, err := clientConn.Read(receivedFromClient)
	if err != nil {
		return
	}
	if tunnelConn != nil {
		tunnelMsg := receivedFromClient[:n]
		log.Println(string(tunnelMsg))
		_, err := tunnelConn.Write(tunnelMsg)
		if err != nil {
			log.Fatalf("sending client request to tunnel: %s\n", err)
		}

		<-clientIsResponded
	}
}

func receiver() {
	received := make([]byte, 1024)
	for {
		n, err := tunnelConn.Read(received)
		if err != nil {
			log.Fatalf("reading response from tunnel: %s\n", err)
		}
		if n == 0 || string(received[:n]) == "HB" {
			continue
		}

		_, err = clientConn.Write(received[:n])
		if err != nil {
			clientConn = nil
		}
		clientIsResponded <- struct{}{}
		continue
	}
}
