package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Start server...")

	// listen on port 8000
	ln, _ := net.Listen("tcp", ":9300")

	// accept connection
	conn, _ := ln.Accept()

	defer conn.Close()
	// run loop forever (or until ctrl-c)
    
	for {
        buf := bufio.NewReader(conn)
		message, _ := buf.ReadString('@')
        c,_ := buf.ReadByte()
		fmt.Println("Received:", string(message))
        fmt.Println("received:", string(c))

		_, err := conn.Write([]byte("recv: " + message + "@"))
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", "recv: "+message)
	}
}
