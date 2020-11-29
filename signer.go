package main

import (
	"fmt"
	"time"
	"strconv"
	//"sort"
	//"strings"
	//"sync"
)

var ExecutePipeline = func(jobs ...job) {
	//wg := &sync.WaitGroup{}
	var in = make(chan interface{}, 10)
	var out = make(chan interface{}, 10)	
	for _, task := range jobs {
		//wg.Add(1)
		go func(task job, in, out chan interface{}) {
			//defer wg.Done()
			task(in, out)
			for message := range out {
				in <- message
			}
		}(task, in, out)
		// task(in, out)
		// for message := range out {
		// 	in <- message
		// }
		time.Sleep(time.Millisecond)
	}
	//fmt.Scanln()
	//wg.Wait()
	time.Sleep(time.Millisecond * 5000)
	close(out)
	close(in)
	
	fmt.Println("ExecutePipeline finished")
}

// var ExecutePipelineSync = func(inputData ...int) {
// 	var results []string
// 	for _, value := range inputData {
// 		singleHash := SingleHash(value)
// 		fmt.Println("SingleHash " + singleHash)
// 		multiHash := MultiHash(singleHash)
// 		fmt.Println("MultiHash " + multiHash)
// 		fmt.Println("-------------------------")
// 		results = append(results, multiHash)
// 	}
// 	totalResult := CombineResults(results)
// 	fmt.Println("CombineResults " + totalResult)
// 	fmt.Println("ExecutePipelineSync finished")
// }

// var SingleHash = func(data int) string {
// 	fmt.Printf("SingleHash input %d\n", data)
// 	strVal := strconv.Itoa(data)
// 	md5Val := DataSignerMd5(strVal)
// 	fmt.Println("SingleHash md5 " + md5Val)
// 	md5Crc32Val := DataSignerCrc32(md5Val)
// 	fmt.Println("SingleHash crc32(md5(data)) " + md5Crc32Val)
// 	crc32Val := DataSignerCrc32(strVal)
// 	fmt.Println("SingleHash crc32(data) " + crc32Val)
// 	return crc32Val + "~" + md5Crc32Val 
// }

// var MultiHash = func(data string) string {
// 	multiHash := ""
// 	fmt.Printf("MultiHash input %s\n", data)
// 	for indx := 0; indx <= 5; indx++ {
// 		indxStrVal := strconv.Itoa(indx)
// 		crc32Val := DataSignerCrc32(indxStrVal + data)
// 		fmt.Println("MultiHash crc32(th + data) " + indxStrVal + " " + crc32Val)
// 		multiHash += crc32Val
// 	}
// 	return multiHash
// }

// var CombineResults = func(results []string) string {
// 	fmt.Printf("CombineResults input %v\n", results)
// 	sort.Slice(results, func(i, j int) bool {
// 		return results[i] < results[j]
// 	})
// 	return strings.Join(results, "_")
// }

var SingleHash = func(in, out chan interface{}) {
	//for val := range in {
		val := <-in
		inputVal := val.(int)
		fmt.Printf("SingleHash input %d\n", inputVal)
		strVal := strconv.Itoa(inputVal)
	 	md5Val := DataSignerMd5(strVal)
		fmt.Println("SingleHash md5 " + md5Val)
		out <- md5Val
		//time.Sleep(time.Millisecond * 100)
	//} 
	//close(out)
	//fmt.Println("SingleHash called")
	//time.Sleep(time.Millisecond * 100)
}

var MultiHash = func(in, out chan interface{}) {
	fmt.Println("MultiHash called")
	//for val := range in {
		val := <-in
		inputVal := val.(string)
		fmt.Printf("MultiHash input %s\n", inputVal)
		//strVal := strconv.Itoa(inputVal)
		//md5Val := DataSignerMd5(strVal)
		result := inputVal + " Finished"
		fmt.Println("MultiHash result " + result)
		out <- result
		//time.Sleep(time.Millisecond * 100)
	//}
}

var CombineResults = func(in, out chan interface{}) {
	fmt.Println("CombineResults called")
}

func main() {
	testResult := "NOT_SET"

	inputData := []int{0}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		//job(CombineResults),
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
	//ExecutePipelineSync(inputData...)
	fmt.Println(testResult)
}