package main

import (
	"fmt"
	"time"
)

func main() {

	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
			fmt.Println(time.Now().Local().Format("2006-01-02"))
		}
	}()

	time.Sleep(time.Millisecond * 3100)
	// ticker.Stop()
	fmt.Println("Ticker stopped")
}
