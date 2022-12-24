package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req, err := http.Get("http://localhost/v1/ping")
			if err != nil {
				fmt.Printf("初始化http客户端处错误:%v", err)
				return
			}
			nByte, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("读取http数据失败:%v", err)
				return
			}
			fmt.Printf("[%d]接收到到值:%v\n", i, string(nByte))
		}(i)
	}

	wg.Wait()
	fmt.Println("请求完毕")
}
