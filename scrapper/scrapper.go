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
// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"
// )

// type requestResult struct {
// 	url     string
// 	status  string
// 	errCode int
// }

// func main() {
// 	ch := make(chan requestResult)
// 	// var results = make(map[string]string) // map으로 처리한경우
// 	var results = make(map[string][]string)

// 	urls := []string{
// 		"https://academy.nomadcoders.co/",
// 		"https://www.airbnb.com/",
// 		"https://www.google.com/",
// 		"https://www.amazon.com/",
// 		"https://www.reddit.com/",
// 		"https://soundcloud.com/",
// 		"https://www.facebook.com/",
// 		"https://www.instagram.com/",
// 	}

// 	for _, url := range urls {
// 		go hitURL(url, ch)
// 	}

// 	for i := 0; i < len(urls); i++ {
// 		result := <-ch
// 		// results[result.url] = result.status // map으로 처리한경우
// 		results[result.url] = append(results[result.url], result.status, strconv.Itoa(result.errCode))
// 	}

// 	for url, status := range results {
// 		// fmt.Println(url, status) // map으로 처리한경우
// 		fmt.Println(url, status[0], status[1])
// 	}
// }

// func hitURL(url string, ch chan<- requestResult) {
// 	resp, err := http.Get(url)
// 	status := "OK"
// 	errCode := 0

// 	if err != nil || resp.StatusCode >= 400 {
// 		status = "FAILED"
// 		errCode = resp.StatusCode
// 	} else {
// 		status = "OK"
// 		errCode = 0
// 	}
// 	ch <- requestResult{url: url, status: status, errCode: errCode}
// }

// go get github.com/PuerkitoBio/goquery
// indeed

package scrapper

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

//https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=python&limit=50

func Scrape(term string) {

	var baseURL string = "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedJob
	ch := make(chan []extractedJob)
	totalPages := getPages(baseURL)
	fmt.Println(totalPages)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, ch)
		// jobs = append(jobs, job...)
	}

	for i := 0; i < totalPages; i++ {
		job := <-ch
		jobs = append(jobs, job...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extraced", len(jobs))
}

func getPage(page int, url string, mainCh chan<- []extractedJob) {

	var jobs []extractedJob
	ch := make(chan extractedJob)
	pageURL := url + "&start=" + strconv.Itoa(page*50)
	fmt.Println("requesting:", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	//slider_item
	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, ch)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-ch
		jobs = append(jobs, job)
	}

	mainCh <- jobs

}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	file.Write(utf8bom)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)
	//https://kr.indeed.com/viewjob?jk=b376351fcbd1750e
	for _, job := range jobs {
		jobSlice := []string{
			"https://kr.indeed.com/viewjob?jk=" + job.id,
			job.title,
			job.location,
			job.salary,
			job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)

	}

}

func extractJob(card *goquery.Selection, ch chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := CleanString(card.Find(".jobTitle>span").Text())
	location := CleanString(card.Find(".companyLocation").Text())
	salary := CleanString(card.Find(".salary-snippet").Text())
	summary := CleanString(card.Find(".job-snippet").Text())

	ch <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}

}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	// fmt.Println(doc)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(s.Html())
		pages = s.Find("a").Length()
	})

	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode, res.Status)
	}
}
