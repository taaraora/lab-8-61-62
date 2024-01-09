package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp4", "127.0.0.1:7890")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	//laddr := conn.LocalAddr().String()
	raddr := conn.RemoteAddr().String()
	defer conn.Close()
	for {
		// Read client request
		msg := make([]byte, 1024)
		n, err := conn.Read(msg)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("client: %s, Received: %s\n", raddr, msg[:n])

		_, err = conn.Write(msg[:n])
		if err != nil {
			fmt.Println("cannot write message back to the client:", err)
			return
		}
	}
}
