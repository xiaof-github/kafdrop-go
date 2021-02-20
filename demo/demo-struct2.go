package main

import (
	"fmt"
)

type cook struct{
	spoon string
	knife string
}

type MyString string

type cookk struct{
	spoon MyString
	knife MyString
}

func main(){
	var cook1 cook = cook{spoon:"so1", knife: "小泉",}
	var cook2 cook = cook{spoon:"so2", knife: "王麻子",}
	fmt.Println(cook1)
	cook1 = cook2
	fmt.Println(cook1)

	var cookk1 cookk = cookk{spoon:"ok1", knife: "小z",}
	var cookk2 cookk = cookk{spoon:"ok2", knife: "王m",}
	fmt.Println(cookk1)
	cookk1 = cookk2
	fmt.Println(cookk1)
}