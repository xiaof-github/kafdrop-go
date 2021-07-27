package main

import (
    "fmt"
    "log"

    req "github.com/imroc/req"
)

func main() {
    reqServer()
}

func reqServer() {
    header := req.Header{
        "Accept": "application/json",
    }
    param := req.Param{
        "key":     "d0bc8a91e590ad02cc518b2968b4c85d",
        "address": "百官街道珠峰新村西一区60幢1401室",
        "city":    "绍兴",
    }
    // 只有url必选，其它参数都是可选
    r, err := req.Get("https://restapi.amap.com/v3/geocode/geo", header, param)
    if err != nil {
        log.Fatal(err)
    }
    v, err := r.ToString()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%v\n", r) // 打印详细信息
    fmt.Println(v)

    resp := r.Response()
    fmt.Println(resp.StatusCode)

}
