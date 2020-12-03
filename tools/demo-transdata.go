package demo

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	ch := make(chan string, 10)
	sendString := func (v string) {
		ch <- v
	}
	recvString := func () string {
		v := <-ch
		fmt.Println("recv:", v)
		return v
	}
	for i:=0; i<10; i++ {
		go sendString(strconv.Itoa(i))
	}
	
	// recv consume data
	for i:=0; i<10; i++ {
		go recvString()		
	}	

    time.Sleep(time.Second)
}