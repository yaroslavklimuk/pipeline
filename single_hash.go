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
		notifChan := make(chan struct{})

		go func(val string, resultArr []string, notifChan chan struct{}) {
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
			notifChan <- struct{}{}
		}(val.(string), results[:], notifChan)

		go func(val string, resultArr []string, notifChan chan struct{}) {
			resultArr[0] = DataSignerCrc32(val)
			notifChan <- struct{}{}
		}(val.(string), results[:], notifChan)

		<- notifChan
		<- notifChan
		close(notifChan)
		out <- strings.Join(results[:], "~")
	}
	close(out)
}
