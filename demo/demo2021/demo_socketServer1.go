package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Start server...")

	// listen on port 8000
	ln, _ := net.Listen("tcp", ":8000")

	// accept connection
	conn, _ := ln.Accept()

	defer conn.Close()
	// run loop forever (or until ctrl-c)
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Received:", string(message))

		_, err := conn.Write([]byte("recv: " + message))
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", "recv: "+message)
	}
}
