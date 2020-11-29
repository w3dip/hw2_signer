package main

import (
	"fmt"
	//"time"
	"sort"
	"strconv"
	"strings"
	"sync"
	"runtime"
)

var ExecutePipeline = func(jobs ...job) {
	runtime.GOMAXPROCS(8)
	wg := &sync.WaitGroup{}
	var in = make(chan interface{}, 100)
	for _, task := range jobs {
		var out = make(chan interface{}, 100)
		wg.Add(1)
		go func(wg *sync.WaitGroup, in, out chan interface{}, task job) {
			defer wg.Done()
			defer close(out)
			task(in, out)
		}(wg, in, out, task)
		in = out
	}
	wg.Wait()
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

var SingleHash = func(in, out chan interface{}) {
	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for data := range in {
		wg.Add(1)
		go func(in interface{}, out chan interface{}, wg *sync.WaitGroup, mutex *sync.Mutex) {
			defer wg.Done()
			data := strconv.Itoa(in.(int))
			crc32Chan := make(chan string)
			go func(data string, out chan string) {
				out <- DataSignerCrc32(data)
				runtime.Gosched()
			}(data, crc32Chan)
			crc32Val := <-crc32Chan
			mutex.Lock()
			md5Val := DataSignerMd5(data)
			mutex.Unlock()
			md5Crc32Val := DataSignerCrc32(md5Val)
			result := crc32Val + "~" + md5Crc32Val
			out <- result
		}(data, out, wg, mutex)
	}
	wg.Wait()
}

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

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(data string, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			results := make([]string, 6)
			mutex := &sync.Mutex{}
			wgItem := &sync.WaitGroup{}
			for i := 0; i <= 5; i++ {
				wgItem.Add(1)
				dataStrVal := strconv.Itoa(i)
				data := dataStrVal + data
				go func(results []string, index int, data string, wgItem *sync.WaitGroup, mutex *sync.Mutex) {
					defer wgItem.Done()
					data = DataSignerCrc32(data)
					runtime.Gosched()
					mutex.Lock()
					results[index] = data
					mutex.Unlock()
				}(results, i, data, wgItem, mutex)
			}
			wgItem.Wait()
			totalResult := strings.Join(results, "")
			out <- totalResult
		}(data.(string), out, wg)
	}
	wg.Wait()
}

// var CombineResults = func(results []string) string {
// 	fmt.Printf("CombineResults input %v\n", results)
// 	sort.Slice(results, func(i, j int) bool {
// 		return results[i] < results[j]
// 	})
// 	return strings.Join(results, "_")
// }

func CombineResults(in, out chan interface{}) {
	var result []string
	for data := range in {
		result = append(result, data.(string))
	}
	sort.Strings(result)
	totalResult := strings.Join(result, "_")
	out <- totalResult
}

func main() {
	testResult := "NOT_SET"

	inputData := []int{0, 1}

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
	//ExecutePipelineSync(inputData...)
	fmt.Println(testResult)
}
