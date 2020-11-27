package demo

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func hello() {
	defer wg.Done()
	fmt.Println("hello")
}

func main() {
	wg.Add(1)
	go hello()
	fmt.Println("main")
	wg.Wait()
}