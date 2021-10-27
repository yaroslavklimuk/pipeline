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
			close(in)
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			for dataRaw := range in {
				data, ok := dataRaw.(string)
				if ok {
					fmt.Print(data)
				}
			}
		}),
	}
	return jobsStack
}

func ExecutePipeline(jobsList ...job) {
	jobsCount := len(jobsList)
	channels := make([]chan interface{}, jobsCount+1)
	canMD5.Store(true)
	channels[jobsCount] = make(chan interface{})
	for i := jobsCount - 1; i >= 0; i-- {
		start := time.Now()
		channels[i] = make(chan interface{})
		wg.Add(1)
		go func(ind int) {
			jobsList[ind](channels[ind], channels[ind+1])
			close(channels[ind+1])
			wg.Done()
		}(i)
		end := time.Since(start)
		fmt.Printf("Started at: %s\n", start)
		fmt.Printf("Finished in: %s\n", end)
	}
	wg.Wait()
}

func main() {
	jobs := MakeJobsList()
	start := time.Now()
	ExecutePipeline(jobs...)
	end := time.Since(start)
	fmt.Printf("Started at: %s\n", start)
	fmt.Printf("Finished in: %s\n", end)
}
