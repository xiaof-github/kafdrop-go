package demo

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
	var atr []*Hello
	for i := 0; i< 10; i++ {
		a := new(Hello)
		atr = append(atr, a)
	}
	time.Sleep(time.Second)
	fmt.Println("main")
}