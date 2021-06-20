package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(params ...job) {
	wg := &sync.WaitGroup{}
	in, out := make(chan interface{}), make(chan interface{})
	for _, j := range params{
		wg.Add(1)
		go func(wg *sync.WaitGroup, j job, in, out chan interface{}) {
			defer wg.Done()
			j(in, out)
			close(out)
		}(wg,j,in,out)
		in = out
		out = make(chan interface{})
	}
	wg.Wait()
	close(out)
}

func SingleHash(in, out chan interface{}){
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for data := range in {
		value, ok := data.(string)
		if !ok {
			value = strconv.Itoa(data.(int))
		}
		wg.Add(1)

		go func(data string) {
			defer wg.Done()
			dh1, dh2 := make(chan string), make(chan string)

			go func() {
				dh1 <- DataSignerCrc32(data)
			}()
			go func() {
				mu.Lock()
				md5 := 	DataSignerMd5(data)
				mu.Unlock()
				dh2 <- DataSignerCrc32(md5)
			}()
			out <- <-dh1 + "~" + <-dh2
		}(value)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}){
	wg := &sync.WaitGroup{}
	for data := range in {
		value, ok := data.(string)
		if !ok {
			value = strconv.Itoa(data.(int))
		}
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			workers := &sync.WaitGroup{}
			hashes := make([]string, 6)
			for th := 0; th < 6.; th++ {
				workers.Add(1)
				go func(th int) {
					defer workers.Done()
					hashes[th] = DataSignerCrc32(strconv.Itoa(th)+data)
				}(th)
			}
			workers.Wait()
			out <- strings.Join(hashes,"")
		}(value)
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}){
	results := make([]string,0,5)
	for data := range in {
		results = append(results, data.(string))
	}
	sort.Strings(results)
	result := strings.Join(results, "_")
	out <- result
}
