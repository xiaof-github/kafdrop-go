package maich

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func a(ch chan string) {
	var c [100]string = [100]string{}
	defer wg.Done()
    for i := 1; i < 10; i++ {
		c[i] = strconv.Itoa(i)
		ch <- c[i]
	}	
}

func b(ch chan string) {
    defer wg.Done()
	for num := range ch {
		fmt.Printf("num:%s, channel len:%d\n",num, len(ch))
		if len(ch) <= 0 { // 如果现有数据量为0，跳出循环
			break
		   //close(ch)    //及时关闭也可以避免死锁
		}
	
	}
	
}

func main() {
	var ch1 chan string
	// var ch2 chan string
	runtime.GOMAXPROCS(2)
	
	ch1 = make(chan string, 100)
	// ch2 = make(chan string)
	wg.Add(1)
	go a(ch1)
	wg.Add(1)
	go b(ch1)
	fmt.Println("main")
    wg.Wait()
}