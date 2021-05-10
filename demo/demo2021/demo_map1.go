package main
 
import (
	"fmt"
)
 
type te struct {
	name string
	val  int
}
func main() {
	x := make(map[int]te)
	for i := 0; i < 3; i++ {
		x[i] = te{name: "a", val:i}
	}
	for k, v := range x {
		v.name = "b"
		v.val += 1
		fmt.Println(k, v)
	}
	fmt.Println("------------change in range-----------")
	for k, v := range x {
		fmt.Println(k, v)
	}

	fmt.Println("------------copy in range-----------")
	for k, v := range x {
		v.name = "b"
		v.val += 1
		x[k] = v
	}

	for k, v := range x {
		fmt.Println(k, v)
	}

}