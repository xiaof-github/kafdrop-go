package agent

import (
	"fmt"

	"github.com/shirou/gopsutil/process"
)

// https://pkg.go.dev/github.com/shirou/gopsutil@v3.21.3+incompatible/process

func main() {

	processes, _ := process.Processes()
	for _, p := range processes {
		fmt.Println(p)
		fmt.Println(p.Name())
	}

}
