package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/host"
)

func main() {

	fmt.Println("day1: ", time.Now().Local().Format("2006-01-02"))
	fmt.Println("day2: ", time.Now().Local().Format("20060102"))
	
	dayBegin := time.Now().Local().Format("2006-01-02") + " 00:00:00"
	// 此格式为固定格式
	// Parse parses a formatted string and returns the time value it represents. The layout defines the format by showing how the reference time, defined to be Mon Jan 2 15:04:05 -0700 MST 2006 would be interpreted if it were the value; it serves as an example of the input format. The same interpretation will then be made to the input string.

	format := "2006-01-02 15:04:05"

	
	timestamp, _ := host.BootTime()
	timeN := time.Now()
	timeNs := time.Now().Unix()

	fmt.Println("boot, now, stamp, dayBegin, ", timestamp, timeN, timeNs, dayBegin)

	stamp, err := time.ParseInLocation(format, dayBegin, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	if (err != nil ){
		fmt.Println("err: ", err)
	}
	fmt.Println("time parse: ", stamp, stamp.Unix())

}