package main

import (
	"fmt"
	"strconv"
)
func main() {
	x := make(map[int]string, 20)
	for i := 0; i < 20; i++ {
		x[i] = strconv.Itoa(i)
	}
    delete(x, 1)
    delete(x, 0)
    fmt.Println(x)
}