package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

const (
	combineMultiHashSeparator string = "_"
	multiHashSeparator        string = ""
	crc32md5Separator         string = "~"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})
	for _, job := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go worker(job, in, out, wg)
		in = out
	}

	wg.Wait()
}

func worker(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	job(in, out)
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	md5QuotaCh := make(chan int, 1)
	for val := range in {
		crc32md5Ch := make(chan interface{}, 1)
		wg.Add(1)
		go func(val string, out chan<- interface{}, quotaCh chan int) {
			defer wg.Done()
			quotaCh <- 1
			md5 := DataSignerMd5(val)
			<-quotaCh

			crc32md5 := DataSignerCrc32(md5)
			crc32md5Ch <- crc32md5
		}(fmt.Sprint(val), crc32md5Ch, md5QuotaCh)

		wg.Add(1)
		go func(val string, in <-chan interface{}, out chan<- interface{}) {
			defer wg.Done()
			crc32Hash := DataSignerCrc32(val)
			crcMd5Hash := fmt.Sprint(<-in)
			out <- crc32Hash + crc32md5Separator + crcMd5Hash
		}(fmt.Sprint(val), crc32md5Ch, out)
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for val := range in {
		wgMultiHashes := &sync.WaitGroup{}
		resulstHashes := make([]string, 6)
		for i := 0; i < 6; i++ {
			wgMultiHashes.Add(1)
			go func(val string, th int) {
				defer wgMultiHashes.Done()

				res := DataSignerCrc32(fmt.Sprint(th) + val)

				mu.Lock()
				resulstHashes[th] = res
				mu.Unlock()
			}(fmt.Sprint(val), i)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			wgMultiHashes.Wait()

			multiHash := strings.Join(resulstHashes, multiHashSeparator)

			out <- multiHash
		}()
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	resultHashes := make([]string, 0, 2)
	for hash := range in {
		resultHashes = append(resultHashes, fmt.Sprint(hash))
	}

	sort.Strings(resultHashes)

	out <- strings.Join(resultHashes, combineMultiHashSeparator)
}
