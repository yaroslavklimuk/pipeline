package main

import (
	"sort"
)

var CombineResults = func(in, out chan interface{}) {
	results := make([]string, 10)
	for val := range in {
		results = append(results, val.(string))
	}
	sort.Strings(results)
	for i := 0; i < len(results); i++ {
		out <- results[i]
	}
	close(out)
}
