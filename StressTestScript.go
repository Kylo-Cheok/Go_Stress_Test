package main

/**
 * @Author: kylo_cheok
 * @Email:  maggic0816@gmail.com
 * @Date:   2022/1/23 20:29
 * @Desc:   Grace under pressure
 */
import (
	"GoProject/StressTest/structs"
	"GoProject/StressTest/tools"
	"fmt"
	"time"
)

func main() {
	num := 100   //协程数
	count := 100 //每个协程请求数

	r := &structs.Request{
		URL:    "https://www.baidu.com",
		Method: "GET",
		Headers: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		Body: "",
	}
	begin := time.Now()
	tools.StartTest(num, count, r)
	fmt.Println("压测总耗时:", time.Since(begin))
}
