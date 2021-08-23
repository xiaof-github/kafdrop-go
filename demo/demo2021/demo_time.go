package main

import (
    "time"
    "fmt"
    "github.com/shirou/gopsutil/host"
    
)

func main(){
    timestamp, _ := host.BootTime()
    fmt.Println(timestamp)
    t,_ := getDayBeginStamp()

    pcRunTime := int((time.Now().Unix() - t) / 60)
    fmt.Println("pc run time:", pcRunTime)
    
}

func getDayBeginStamp() (int64, error) {
    // 时间模板格式的年月日为固定值，不可修改
    format := "2006-01-02 15:04:05"
    dayBegin := time.Now().Local().Format("2006-01-02") + " 00:00:00"
    stamp, err := time.ParseInLocation(format, dayBegin, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
    if err != nil {
        fmt.Println("err: ", err)
    }
    return stamp.Unix(), err
}
