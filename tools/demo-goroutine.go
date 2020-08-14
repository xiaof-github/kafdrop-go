package maing

import (
	"fmt"
	"runtime"
	"time"
)

func a() {
    for i := 1; i < 10; i++ {
        fmt.Println("A:", i)
    }
}

func b() int {
    for i := 1; i < 10; i++ {
        fmt.Println("B:", i)
    }
    return 0
}

func main() {
    runtime.GOMAXPROCS(2)
    go a()
    go b()
    time.Sleep(time.Second)
}