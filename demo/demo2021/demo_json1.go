package main

import (	
	"fmt"
	"encoding/json"
)

func main() {
	map2json()
}

func map2json(){
	m := map[string]string{"name": "张三", "id": "001"}
	mjson,_ := json.Marshal(m)
	fmt.Println(mjson)
	mString := string(mjson)
	fmt.Println("mString: ", mString)
}