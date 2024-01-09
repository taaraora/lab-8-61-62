package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"strconv"
)

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Printf("args: %v\n", args)
	if len(args) != 2 {
		fmt.Println("Usage: echoclient <server-host> <server-port>")
		os.Exit(1)
	}

	a, err := netip.ParseAddr(args[0])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if !a.IsValid() {
		fmt.Println("Error: invalid address")
		os.Exit(1)
	}

	p, err := strconv.ParseUint(args[1], 10, 32)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if p > 65535 {
		fmt.Println("Error: invalid port")
		os.Exit(1)
	}

	addr := args[0] + ":" + args[1]
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// sending goroutine
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := scanner.Bytes()
			_, err = conn.Write(msg)
			if err != nil {
				log.Fatal(err)
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()

	// receiving goroutine
	go func() {
		for {
			msg := make([]byte, 1024)
			n, err := conn.Read(msg)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Printf("Received: %s\n", msg[:n])
		}
	}()

	select {}

}
