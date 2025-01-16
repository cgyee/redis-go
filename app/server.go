package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit
var s *Store

type Conn struct {
	conn net.Conn
}

type Cmd struct {
	name string
	raw  []byte
	data []byte
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	s = GetInstance()

	l, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	ch := make(chan net.Conn, 1000)
	for i := len(os.Args) - 1; i > 0; i-- {
		if os.Args[i] == "--dir" {
			s.Set("dir", "$", []byte(os.Args[i+1]))
		} else if os.Args[i] == "--dbfilename" {
			s.Set("dbfilename", "$", []byte(os.Args[i+1]))
		}
	}
	r := s.Get("dir")
	fmt.Println("store main: ", string(r.Value))
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
			}
		default:
			continue
		}
	}
}

func handleRequest(conn net.Conn) {
	buff := make([]byte, 4096)
	for {
		l, err := conn.Read(buff)
		if err != nil {
			if err.Error() == "EOF" {
				// fmt.Println("Connection closed")
				continue
			}
			// fmt.Println("Error reading:", err.Error())
			continue
		} else if l == 0 {
			fmt.Println("No data received")
		} else {
			_, msg := Handle(buff[:l])
			conn.Write(msg)
		}
	}
}
