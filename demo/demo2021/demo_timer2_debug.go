package main

import (
    "fmt"
    "time"

    "reflect"

    log "github.com/sirupsen/logrus"
)

func main() {

    var queue = make(chan int, 5)

    ticker2 := time.NewTicker(time.Second * time.Duration(10))
    defer ticker2.Stop()

    // cct := <-queue
    // fmt.Println("cct type: ", reflect.TypeOf(cct))
    var tr int
    fmt.Println("t type: ", reflect.TypeOf(tr))

    gpro := 100
    go func(ti *time.Ticker) {
        for t := range ti.C {
            log.Info("Statistics Tick send at ", t)

            // 查看进程间隔到达时，累加一次进程运行时间，时间间隔为60s的倍数
            gpro++
            fmt.Println(len(queue))
            if len(queue) >= 4 {
                var cct int
                cct = <-queue
                fmt.Println("cct : ", cct)
            }
            queue <- gpro
            fmt.Println("send gpro", t)
        }
    }(ticker2)

    ticker3 := time.NewTicker(time.Minute * time.Duration(1))
    defer ticker3.Stop()

    // 定时拉取配置
    go func(ti *time.Ticker) {
        fmt.Println("ticker get conf")
        for t := range ti.C {
            log.Info("Get conf Tick send at ", t)
            for i := 0; i < 100; i++ {
                c := i + 1
                if c%99 == 0 {
                    cc := <-queue
                    fmt.Println(cc)
                }
            }
        }
    }(ticker3)

    //  test
    for {
        time.Sleep(time.Second)
    }
}
