package main

import (
	"bufio"
	"fmt"
	"net"
)

var s string = "abcde"

func main() {
	connect, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(connect, s+"\n")
	message, err := bufio.NewReader(connect).ReadString('\n')
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(message)
}
