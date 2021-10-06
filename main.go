/*
// URL checker
package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("request failed")

func main() {
	var results = make(map[string]string)

	urls := []string{
		"https://academy.nomadcoders.co/",
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		result := "OK"
		err := hitURL(url)
		if err != nil {
			result = "FAILED"
		}
		results[url] = result
	}
	for i, v := range results {
		fmt.Println(v, i)
	}
}

func hitURL(url string) error {
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		fmt.Println(url, err, resp.StatusCode)
		return errRequestFailed
	}
	return nil
}
*/

// // goroutine
// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// var wg sync.WaitGroup

// func main() {
// 	go sexyCount("AAA")
// 	fmt.Println("******************")
// 	go sexyCount("BBB ")

// 	wg.Add(2)
// 	wg.Wait()
// }

// func sexyCount(person string) {
// 	for i := 0; i < 5; i++ {
// 		fmt.Println(person, "is sexy", i)
// 		time.Sleep(time.Millisecond * 1000)
// 	}
// 	wg.Done()
// }

// // channel
// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	c := make(chan string)
// 	people := []string{"AAA", "BB ", "CCC", "DDD", "E  "}
// 	for _, person := range people {
// 		fmt.Println("sending data:", person)
// 		go isSexy(person, c)
// 	}

// 	for i := 0; i < len(people); i++ {
// 		// fmt.Println("Waiting No:", i)
// 		fmt.Println("get from channel:", <-c) // 채널은 자동으로 wait 상태로 기다림, blocking operation
// 	}
// }

// func isSexy(person string, channel chan string) {
// 	time.Sleep(time.Second)
// 	channel <- person + " is sexy"
// }

// URL checker
package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type requestResult struct {
	url     string
	status  string
	errCode int
}

func main() {
	ch := make(chan requestResult)
	// var results = make(map[string]string) // map으로 처리한경우
	var results = make(map[string][]string)

	urls := []string{
		"https://academy.nomadcoders.co/",
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.amazon.com/",
		"https://www.reddit.com/",
		"https://soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
	}

	for _, url := range urls {
		go hitURL(url, ch)
	}

	for i := 0; i < len(urls); i++ {
		result := <-ch
		// results[result.url] = result.status // map으로 처리한경우
		results[result.url] = append(results[result.url], result.status, strconv.Itoa(result.errCode))
	}

	for url, status := range results {
		// fmt.Println(url, status) // map으로 처리한경우
		fmt.Println(url, status[0], status[1])
	}
}

func hitURL(url string, ch chan<- requestResult) {
	resp, err := http.Get(url)
	status := "OK"
	errCode := 0

	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
		errCode = resp.StatusCode
	} else {
		status = "OK"
		errCode = 0
	}
	ch <- requestResult{url: url, status: status, errCode: errCode}
}
