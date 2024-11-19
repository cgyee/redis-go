package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		if conn != nil {
			handleRequest(conn)
		}
	}
}

func handleRequest(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		dataLength, err := conn.Read(buff)
		defer conn.Close()

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed")
			}
			fmt.Println("Error reading:", err.Error())
		}
		if dataLength == 0 {
			fmt.Println("No data received")
		} else {

			msg := strings.Split(string(buff), "\r\n")
			fmt.Printf("Message: %v", msg)
			fmt.Println(msg)

			for _, m := range msg {
				switch strings.TrimSpace(m) {
				case "PING":
					conn.Write([]byte("+PONG\r\n"))
				default:
					fmt.Println("Received data: ", string(m))
				}
			}
		}
	}
}
