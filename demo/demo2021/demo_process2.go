package main

import (
    "fmt"
    "time"
    "strings"

    "github.com/shirou/gopsutil/process"
    "github.com/shirou/gopsutil/host"
)

// https://pkg.go.dev/github.com/shirou/gopsutil@v3.21.3+incompatible/process

func main() {

    processes, _ := process.Processes()
    for _, p := range processes {
        name, _ := p.Name()
        createTime, err := p.CreateTime()
        if (err != nil) {
            fmt.Println("err: ", name, err, createTime)
        }
        runTime := (time.Now().Unix() -createTime/1000)/60
        fmt.Println(p, name, createTime/1000, runTime)

    }

    timestamp, _ := host.BootTime()
    t := time.Unix(int64(timestamp), 0)
    fmt.Println(t.Local().Format("2006-01-02 15:04:05"))
    s := time.Now().Local().Format("2006-01-02")
    fmt.Println(strings.ReplaceAll(s, "-", ""))
    fmt.Println("runtime: ", (time.Now().Unix() - int64(timestamp))/60)
}
