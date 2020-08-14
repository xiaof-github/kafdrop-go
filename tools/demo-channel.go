package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

func a(ch chan string) {
	var c [100]string = [100]string{}
    for i := 1; i < 10; i++ {
		c[i] = strconv.Itoa(i)
		ch <- c[i]
    }
}

func b(ch chan string) {
    
	x := <- ch
	fmt.Println(x)
}

func main() {
	var ch1 chan string
	// var ch2 chan string
	runtime.GOMAXPROCS(2)
	
	ch1 = make(chan string)
	// ch2 = make(chan string)
    go a(ch1)
    go b(ch1)
    time.Sleep(time.Second)
}