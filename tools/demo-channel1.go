package demo

import "fmt"
func main() {
    c := make(chan struct{})
    close(c)
    select {
    case c <- struct{}{}: // 若此分支被选中，则产生一个恐慌
    case <-c: fmt.Println("recv")
    }
}