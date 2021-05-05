package main

import (
	"fmt"	
)

func main() {

	s := "abc"
	b := s
	s = "fd"
	
	fmt.Println(s)
	fmt.Println(b)
}
