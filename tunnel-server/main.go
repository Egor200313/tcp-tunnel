package main

import (
	"flag"
	"log"
	"net"
	"sync"
)

var (
	clientPort = flag.String("client-port", "9090", "local server port for public clients requests")
	tunnelPort = flag.String("tunnel-port", "9999", "local server port to accept tunnel agent connections")
	connMap    = &sync.Map{}
)

func main() {
	flag.Parse()
	tunnelLoop()
	//for {
	clientLoop()
	//}
}

func clientLoop() {
	l, err := net.Listen("tcp", "0.0.0.0:"+*clientPort)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("waiting for client...")
	//for {
	conn, err := l.Accept()
	if err != nil {
		log.Fatalf("accepting client connection: %s\n", err)
	}
	log.Println("got client connection")
	handleClientConnection(conn)
}

func tunnelLoop() {
	l, err := net.Listen("tcp", "localhost:"+*tunnelPort)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("waiting for agent...")
	//for {
	conn, err := l.Accept()
	if err != nil {
		log.Fatalf("accepting agent connection: %s\n", err)
	}
	log.Println("got agent connection")
	connMap.Store("1", conn)
	//}
	//for {
	//conn.Write([]byte("hb"))
	//}
}

func handleClientConnection(c net.Conn) {
	defer c.Close()

	//for {
	received := make([]byte, 1024)
	n, err := c.Read(received)
	//n, err := io.ReadFull(c, received)
	log.Println(n)
	if err != nil {
		log.Fatalf("reading client request: %s\n", err)
	}
	//log.Println(received)
	if outConn, ok := connMap.Load("1"); ok {
		outMsg := received[:n]
		log.Println(string(outMsg), len(outMsg))
		_, err := outConn.(net.Conn).Write(outMsg)
		if err != nil {
			log.Fatalf("sending client request: %s\n", err)
		}

		received = make([]byte, 1024)
		_, err = outConn.(net.Conn).Read(received)
		if err != nil {
			log.Fatalf("reading server response: %s\n", err)
		}
		c.Write(received)
	}
	//time.Sleep(5 * time.Second)
	//}

}
