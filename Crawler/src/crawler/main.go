package main

import (
	"runtime"
	"sync"
	"fmt"
	"time"
	"strings"
	"log"
	"github.com/temoto/robotstxt"
)



func checkDuplicate(in chan Graph) chan Graph {
	ch := make(chan Graph)
	visited := map[string]int{}

	report := func(gr Graph) {
		if gr.from == gr.to {
			return
		}
		fmt.Println(visited[gr.from], visited[gr.to])
	}

	go func() {
		defer close(ch)
		for {
			select {
			case x := <-in:
				if _, ok := visited[x.to]; !ok {
					visited[x.to] = len(visited)
					ch <- x
				}
				report(x)

			case <- DefaultDone:
				return
			}
		}
	}()

	return ch
}

func checkRobo(url string) bool {
	path := url
	if idx := strings.Index(url[10:], "/"); idx != -1 {
		url = url[:10 + idx]
	}

	RobotsLock.RLock()
	robo, ok := DefaultRobots[url]
	if !ok {
		RobotsLock.RUnlock()
		resp, err := DefaultClient.Get(url + "/robots.txt")
		if err != nil {
			log.Println(err)
			return false
		}
		defer resp.Body.Close()
		rd, err := robotstxt.FromResponse(resp)
		if err != nil {
			log.Println(err)
			return false
		}
		robo = rd.FindGroup(DefaultUserAgent)

		RobotsLock.Lock()
		DefaultRobots[url] = robo
		RobotsLock.Unlock()
	} else {
		RobotsLock.RUnlock()
	}

	return robo.Test(path)
}

func addressUri(url string) {
	if !checkRobo(url) {
		return
	}

	log.Println("visiting:", url, "goroutines:", runtime.NumGoroutine())

	resp, err := DefaultClient.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	for _, link := range parseLinks(resp) {
		input  <- Graph{url, link}
	}
}

func solver() {
	input <- Graph{DefaultInitialWeb, DefaultInitialWeb}
	var wg sync.WaitGroup
	q := smallBuffer(checkDuplicate(input))
	wg.Add(DefaultNumWorkers)
	worker := func() {
		defer wg.Done()
		for {
			select {
			case uri := <-q:
				addressUri(uri.to)
			case <-time.After(DefaultCrawlDelay):
				return
			}
		}
	}
	for i := 0; i < DefaultNumWorkers; i++ {
		go worker()
	}
	wg.Wait()
	close(input)
	close(DefaultDone)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	solver()
}