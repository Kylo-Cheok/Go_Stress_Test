package tools

import (
	"GoProject/StressTest/structs"
	"fmt"
	"sort"
	"sync"
	"time"
)

/**
 * @Author: kylo_cheok
 * @Email:  maggic0816@gmail.com
 * @Date:   2022/2/19 18:18
 * @Desc:   Grace under pressure
 */

func ReceivingResults(concurrent int, ch <-chan *structs.RequestResults, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	var stopChan = make(chan bool)
	// 时间
	var (
		processingTime  float64 // 处理总时间
		maxTime         float64 // 最大时长
		minTime         float64 // 最小时长
		successNum      int     // 成功处理数，code为200
		failureNum      int     // 处理失败数，code不为200
		chanIDLen       int     // 并发数
		chanIDs         = make(map[int]bool)
		mutex           = sync.RWMutex{}
		requestTimeList []float64
	)
	statTime := time.Now().Second()
	// 错误码/错误个数
	// 定时1秒输出一次计算结果
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := time.Now().Second()
				mutex.Lock()
				go calculateData(concurrent, processingTime, endTime-statTime, maxTime, minTime, successNum, failureNum,
					chanIDLen)
				mutex.Unlock()
			case <-stopChan:
				// 处理完成
				return
			}
		}
	}()
	header()
	for data := range ch {
		mutex.Lock()
		processingTime = processingTime + data.RT
		if maxTime <= data.RT {
			maxTime = data.RT
		}
		if minTime == 0 {
			minTime = data.RT
		} else if minTime > data.RT {
			minTime = data.RT
		}
		// 是否请求成功
		if data.Succeed == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}
		if _, ok := chanIDs[data.ID]; !ok {
			chanIDs[data.ID] = true
			chanIDLen = len(chanIDs)
		}
		requestTimeList = append(requestTimeList, data.RT)
		mutex.Unlock()
	}
	// 数据全部接受完成，停止定时输出统计数据
	stopChan <- true
	endTime := time.Now().Second()
	requestTime := endTime - statTime
	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIDLen)

	footer(concurrent, successNum, failureNum, requestTimeList)
}

// calculateData 计算qps 计算1s内 请求的耗时
func calculateData(concurrent int, processingTime float64, requestTime int, maxTime, minTime float64, successNum, failureNum int,
	chanIDLen int) {
	if processingTime == 0 {
		processingTime = 1
	}
	var (
		qps              float64
		averageTime      float64
		requestTimeFloat float64
	)
	// 平均 每个协程成功数*总协程数据/总耗时 (每秒)
	if processingTime != 0 {
		qps = float64(successNum*concurrent*1000) / processingTime
	}
	// 平均时长 总耗时/总请求数
	if successNum != 0 && concurrent != 0 {
		averageTime = processingTime / float64(successNum)
	}
	requestTimeFloat = float64(requestTime)
	// 打印的时长都为毫秒
	table(successNum, failureNum, qps, averageTime, maxTime, minTime, requestTimeFloat, chanIDLen)
}

// header 打印表头信息
func header() {
	// 打印的时长都为毫秒 总请数
	fmt.Println("─────┬───────┬───────┬───────┬────────┬────────┬────────┬────────")
	fmt.Println(" 耗时│ 并发数│ 成功数│ 失败数│   qps  │最长耗时│最短耗时│平均耗时")
	fmt.Println("─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────")
}

// table 打印表格
func table(successNum, failureNum int, qps, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat float64, chanIDLen int) {

	// 打印的时长都为毫秒
	result := fmt.Sprintf("%4.0fs│%7d│%7d│%7d│%8.2f│%8.2f│%8.2f│%8.2f",
		requestTimeFloat, chanIDLen, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime)
	fmt.Println(result)
	return
}

// footer 打印表头信息
func footer(concurrent int, successNum int, failureNum int, requestTimeList []float64) {
	// 打印的时长都为毫秒 总请数
	fmt.Printf("\n\n")
	fmt.Println("*************************  结果 stat  ****************************")
	fmt.Println("处理协程数量:", concurrent, "请求总数:", successNum+failureNum, "成功数:", successNum, "失败数:", failureNum)
	printTop(requestTimeList)
	fmt.Println("*************************  结果 end   ****************************")
	fmt.Println()
}

// printTop 排序后计算 avg top 90 95 99
func printTop(rtList []float64) {
	if rtList == nil {
		return
	}
	sort.Float64s(rtList)
	sum := 0.0
	for _, val := range rtList {
		sum += val
	}
	fmt.Println("avg:", fmt.Sprintf("%.3f", sum/float64(len(rtList))))
	fmt.Println("tp90:", fmt.Sprintf("%.3f", rtList[int(float64(len(rtList))*0.90)]))
	fmt.Println("tp95:", fmt.Sprintf("%.3f", rtList[int(float64(len(rtList))*0.95)]))
	fmt.Println("tp99:", fmt.Sprintf("%.3f", rtList[int(float64(len(rtList))*0.99)]))
}
