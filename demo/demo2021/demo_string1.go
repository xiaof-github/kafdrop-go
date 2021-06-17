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

    t := "23333333daccccccccccccccc"

    fmt.Println(PreStr(t, 10))
}

func PreStr(str string, length int) string {
	var i int
	for i = range str {
		if length == i{
			break
		}		
	}
	return str[:i]
}