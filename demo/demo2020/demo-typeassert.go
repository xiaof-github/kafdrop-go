package main

import "fmt"
type Writer interface {
    Write(buf []byte) (int, error)
}
type DummyWriter struct{}
func (DummyWriter) Write(buf []byte) (int, error) {
    return len(buf), nil
}
func main() {
    var x interface{} = DummyWriter{}
    // y的动态类型为内置类型string。
    var y interface{} = "abc"
    var w Writer
    var ok bool
    // DummyWriter既实现了Writer，也实现了interface{}。
    w, ok = x.(Writer)
	fmt.Println(w, ok) // {} true
	x, ok = w.(interface{})
	fmt.Println(x, ok)
}