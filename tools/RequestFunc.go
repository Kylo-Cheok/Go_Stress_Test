package tools

/**
 * @Author: kylo_cheok
 * @Email:  maggic0816@gmail.com
 * @Date:   2022/2/8 17:00
 * @Desc:   Grace under pressure
 */

import (
	"GoProject/StressTest/structs"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func HttpRequest(request *structs.Request) (float64, string, int) {
	method := request.Method
	url := request.URL
	body := request.GetBody()
	headers := request.Headers

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, "", 0
	}
	client := &http.Client{Timeout: 5 * time.Second}
	for key, header := range headers {
		req.Header.Set(key, header)
	}

	begin := time.Now()
	resp, err := client.Do(req)
	respTime := time.Since(begin).Milliseconds()

	//defer resp.Body.Close()

	resCode := resp.StatusCode
	if resCode != 200 {
		fmt.Println("请求失败:", err)
		return 0, "", 0
	}
	data, err := ioutil.ReadAll(resp.Body)

	return float64(respTime), string(data), resCode
}

func SendHttp(id int, request *structs.Request, ch chan *structs.RequestResults, count int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	var isSucceed bool
	for i := 0; i < count; i++ {
		rt, data, code := HttpRequest(request)
		if code == 200 {
			isSucceed = true
		} else {
			isSucceed = false
		}
		requestResults := &structs.RequestResults{
			ID:       id,
			RT:       rt,
			Succeed:  isSucceed,
			RespData: data,
		}
		ch <- requestResults
	}
}

func StartTest(concurrency int, count int, request *structs.Request) {
	// 设置接收数据缓存
	ch := make(chan *structs.RequestResults, 1000)
	var (
		wg          sync.WaitGroup // 发送数据完成
		wgReceiving sync.WaitGroup // 数据处理完成
	)
	wgReceiving.Add(1)
	go ReceivingResults(concurrency, ch, &wgReceiving)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go SendHttp(i, request, ch, count, &wg)
	}
	// 等待所有的数据都发送完成
	wg.Wait()
	// 延时1毫秒 确保数据都处理完成了
	time.Sleep(1 * time.Millisecond)
	close(ch)
	// 数据全部处理完成了
	wgReceiving.Wait()
	return
}
