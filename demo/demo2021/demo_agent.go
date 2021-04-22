package agent

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

type cpuInfo struct {
	Name          string
	NumberOfCores uint32
	ThreadCount   uint32
	DeviceID      string
}

type Win32_ComputerSystemProduct struct {
	Caption           string
	Description       string
	IdentifyingNumber string
	Name              string
	SKUNumber         string
	Vendor            string
	Version           string
	UUID              string
}

func getCPUInfo() {

	var cpuinfo []cpuInfo

	err := wmi.Query("Select * from Win32_Processor", &cpuinfo)
	if err != nil {
		return
	}
	fmt.Printf("Cpu info =", cpuinfo)
}

func get_ComputerSystemProduct() {
	var product []Win32_ComputerSystemProduct
	err := wmi.Query("Select * from Win32_ComputerSystemProduct", &product)
	if err != nil {
		return
	}
	fmt.Println()
	fmt.Println("product =", product)
}
func main() {
	getCPUInfo()
	get_ComputerSystemProduct()
}
