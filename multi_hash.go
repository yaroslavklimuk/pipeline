package main

import (
	"strconv"
	"strings"
)

var MultiHash = func(in, out chan interface{}) {
	for val := range in {
		var results [6]string
		notifChan := make(chan struct{}, 6)
		for i := 0; i <= 5; i++ {
			go func(item string, ind int, resultArr []string, notChan chan struct{}) {
				resultArr[ind] = strconv.Itoa(ind) + DataSignerCrc32(item)
				notChan <- struct{}{}
			}(val.(string), i, results[:], notifChan)
		}
		for i := 0; i <= 5; i++ {
			<- notifChan
		}
		out <- strings.Join(results[:], "")
	}
	close(out)
}
