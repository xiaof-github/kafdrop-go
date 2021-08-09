package main

import (
    "fmt"
    "sync"
)

func main() {
    var count = 0
    var wg sync.WaitGroup
    var mu sync.Mutex
    //十个协程数量
    n := 10
    wg.Add(n)
    for i := 0; i < n; i++ {
        go func() {
            defer wg.Done()
            //1万叠加
            for j := 0; j < 10000; j++ {
                mu.Lock()
                count++
                mu.Unlock()
            }
        }()
    }
    wg.Wait()
    fmt.Println(count)
}
