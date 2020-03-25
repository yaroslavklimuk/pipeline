package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

type job func(in, out chan interface{})

const (
	MaxInputDataLen = 100
)

func MakeJobsList() []job {
	jobsStack := []job{
		job(func(in, out chan interface{}) {
			inputData := []int{0, 1, 1, 2, 3, 5, 8}
			for _, fibNum := range inputData {
				out <- strconv.Itoa(fibNum)
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}
	return jobsStack
}

func ExecutePipeline(jobsList ...job) {
	jobsCount := len(jobsList)
	channels := make([]chan interface{}, jobsCount+1)
	channels[0] = make(chan interface{})
	canMD5.Store(true)

	i := 0
	channels[i+1] = make(chan interface{})
	func(in, out chan interface{}) {
		inputData := []int{0, 1, 1, 2, 3, 5, 8}
		for _, fibNum := range inputData {
			out <- strconv.Itoa(fibNum)
		}
	}(channels[i], channels[i+1])

	i = i + 1
	channels[i+1] = make(chan interface{})
	SingleHash(channels[i], channels[i+1])

	i = i + 1
	channels[i+1] = make(chan interface{})
	MultiHash(channels[i], channels[i+1])

	i = i + 1
	channels[i+1] = make(chan interface{})
	CombineResults(channels[i], channels[i+1])

	// for i := 0; i <= jobsCount-1; i = i + 1 {
	// 	singleJob := jobsList[jobsCount-1-i]
	// 	channels[i+1] = make(chan interface{})
	// 	singleJob(channels[i], channels[i+1])
	// }

	// for i := 0; i <= jobsCount-1; i = i + 1 {
	// 	singleJob := jobsList[i]
	// 	channels[i+1] = make(chan interface{})
	// 	singleJob(channels[i], channels[i+1])
	// }
	wg.Wait()
	for item := range channels[jobsCount] {
		stv := fmt.Sprintf("%v", item)
		fmt.Println(stv)
	}
}

func main() {
	jobs := MakeJobsList()
	start := time.Now()
	ExecutePipeline(jobs...)
	end := time.Since(start)
	fmt.Printf("Started at: %s\n", start)
	fmt.Printf("Finished in: %s\n", end)
}
