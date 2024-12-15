package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	ch := make(chan net.Conn, 1000)
	go func() {
		for {
			if conn, err := l.Accept(); err == nil {
				ch <- conn
			}
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}
		}
	}()

	for {
		select {
		case conn, ok := <-ch:
			if ok {
				go handleRequest(conn)
			} else {
				continue
			}
		default:
			continue
		}
	}
}

func handleRequest(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		dataLength, err := conn.Read(buff)
		if err != nil {
			if err.Error() == "EOF" {
				// fmt.Println("Connection closed")
				continue
			}
			// fmt.Println("Error reading:", err.Error())
			continue
		} else if dataLength == 0 {
			fmt.Println("No data received")
		} else {
			_, msg := Read(buff[:])
			conn.Write(msg)
		}
	}
}
