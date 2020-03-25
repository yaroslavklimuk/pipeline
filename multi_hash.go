package main

import (
	"strconv"
	"strings"
)

var MultiHash = func(in, out chan interface{}) {
	for val := range in {
		var results [6]string
		notifChan := make(chan int)
		wg.Add(6)
		for i := 0; i <= 5; i++ {
			go func(item string, ind int, resultArr []string, notChan chan int) {
				defer wg.Done()
				resultArr[ind] = strconv.Itoa(ind) + DataSignerCrc32(item)
				notChan <- 1
			}(val.(string), i, results[:], notifChan)
		}

		resultsCount := 0
		for one := range notifChan {
			resultsCount += one
			if resultsCount == 5 {
				out <- strings.Join(results[:], "")
				break
			}
		}
	}
}
