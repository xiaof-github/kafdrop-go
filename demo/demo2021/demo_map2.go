package main

import (
	"fmt"
	"strconv"
)
func main() {
	x := make(map[string]int, 20)
	for i := 0; i < 20; i++ {
		x[strconv.Itoa(i)] = i 
	}
    fmt.Println(x)
    conf := []string{"1","2","3"}
    fmt.Println(conf)
    cutProcess(x, conf)
    fmt.Println("after cut", x)
}

func cutProcess(runTimeMap map[string]int, processConf []string) {
    proMap := make(map[string]int, len(processConf))

    for _, v := range processConf {
        proMap[v] = 0
    }

    for k, _ := range runTimeMap {
        _, ok := proMap[k]
        if !ok {
            delete(runTimeMap, k)
        }
    }
}