package main

import (
	"strings"
	"sync/atomic"
	"time"
)

var canMD5 atomic.Value

var SingleHash = func(in, out chan interface{}) {
	for val := range in {
		var results [2]string
		notifChan := make(chan int)
		wg.Add(2)

		go func(val string, resultArr []string, notifChan chan int) {
			defer wg.Done()
			for {
				if true == canMD5.Load() {
					canMD5.Store(false)
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			md5v := DataSignerMd5(val)
			canMD5.Store(true)
			resultArr[1] = DataSignerCrc32(md5v)
			notifChan <- 1
		}(val.(string), results[:], notifChan)

		go func(val string, resultArr []string, notifChan chan int) {
			defer wg.Done()
			resultArr[0] = DataSignerCrc32(val)
			notifChan <- 1
		}(val.(string), results[:], notifChan)

		resultsCount := 0
		for one := range notifChan {
			resultsCount += one
			if resultsCount == 2 {
				out <- strings.Join(results[:], "~")
				break
			}
		}
	}
}
