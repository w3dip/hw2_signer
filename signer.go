package main

import (
	"fmt"
	"time"
)

var in = make(chan interface{}, 2)
var out = make(chan interface{}, 2)

var ExecutePipeline = func(jobs ...job) {
	for _, task := range jobs {
		task(in, out)
	}
	fmt.Scanln()
	fmt.Println("ExecutePipeline finished")
}

// var SingleHash = func(data int) string {

// }

// var MultiHash = func(data int) string {

// }

var SingleHash = func(in, out chan interface{}) {
	for val := range out {
		inputVal := val.(int)
		fmt.Println("SingleHash input %d", inputVal)
		//md5Val := DataSignerMd5(val.(string))
		//fmt.Println("SingleHash md5 " + md5Val)
		//out <- md5Val
		//time.Sleep(time.Millisecond * 100)
	} 
	fmt.Println("SingleHash called")
	time.Sleep(time.Millisecond * 100)
}

var MultiHash = func(in, out chan interface{}) {
	fmt.Println("MultiHash called")
}

var CombineResults = func(in, out chan interface{}) {
	fmt.Println("CombineResults called")
}

func main() {
	testResult := "NOT_SET"

	inputData := []int{0,1}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			}
			testResult = data
		}),
	}

	ExecutePipeline(hashSignJobs...)
	fmt.Println(testResult)
}