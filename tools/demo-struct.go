package main

import (
	"fmt"
	"time"
)
type Hello struct {
	start int
	end int
}

func hello() {
	var a int
	var b int
	a = 1
	b = 2
	hello := Hello{
		start: a,
		end: b,
	}
	fmt.Println("hello", hello)
}

func main() {

	go hello()
	time.Sleep(time.Second)
	fmt.Println("main")
}