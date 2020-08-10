package mainb

import (
	"fmt"
	"time"
)

func main () {
	timer1 := time.NewTimer(3* time.Second)
	// 执行一次后退出
	timer1.Reset(1 * time.Second)
	fmt.Println(time.Now())
	fmt.Println(<-timer1.C)

	// 定时执行，不退出
	ticker := time.NewTicker(1 * time.Second)
    i := 0
    // 子协程
    go func() {
        for {
            //<-ticker.C
            i++
            fmt.Println(<-ticker.C)
            if i == 5 {
                //停止
                ticker.Stop()
            }
        }
    }()

	for {

	}
}