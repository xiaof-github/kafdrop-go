package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	map2json()
	struct2json()
}

type ProRunTime struct {
	name    string
	runTime int `json:"run_time"`
}

type ProData struct {
	Version string `json:"version"`
	// process      []ProRunTime
	// schoolCode   string `json:"school_code"`
	// buildCode    string `json:"build_code"`
	// houseNum     string `json:"house_num"`
	// uuid         string
	Sn string
	// days         string
	// runTime      int   `json:"run_time"`
	// lastSendTime int64 `json:"last_send_time"`
	// sendTime     int64 `json:"send_time"`
}

func map2json() {
	m := map[string]string{"name": "张三", "id": "001"}
	mjson, _ := json.Marshal(m)
	fmt.Println(mjson)
	mString := string(mjson)
	fmt.Println("mString: ", mString)
}

func struct2json() {
	p := ProData{Version: "V1.0.0", Sn: "11"}
	fmt.Printf("%v\n", p)
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("json error:", err)
	}
	fmt.Println("data", string(data))

}
