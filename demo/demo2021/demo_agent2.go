package agent

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	// "github.com/shirou/gopsutil/mem"  // to use v2
)

func main() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)

	h, _ := host.HostID()
	fmt.Println("host id")
	fmt.Println(h)

	v2, _ := host.Info()
	fmt.Println("host info")
	fmt.Println(v2)
}
