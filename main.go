package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"firefly/fetcher"
	"firefly/safetrie"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"
)

const (
	essayClass        = ".caas-body"
	numOfFetchWorkers = 2
	numOfCountWorkers = 10
)

type FireflyApp struct {
	wordsFilename string
	urlsFilename  string
	wordsBank     *safetrie.SafeTrie
	fetcher       fetcher.Fetcher
}

func isAlphabetic(s string) bool {
	return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(s)
}

func (f *FireflyApp) createDict() {
	f.wordsBank = safetrie.NewSafeTrie()
	file, err := os.Open("./" + f.wordsFilename)
	if err != nil {
		fmt.Println("Error opening words file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) > 2 && isAlphabetic(word) {
			f.wordsBank.Put(strings.ToLower(word), true)
		}
	}
}

func (f *FireflyApp) fetchEssayWorker(urlsChan <-chan string, essaysChan chan<- string, wg *sync.WaitGroup) {
	for url := range urlsChan {
		resp, err := f.fetcher.GetWithRetry(url, 10)
		if err != nil {
			fmt.Println("error getting url: ", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print("error reading resp body: ", err)
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			fmt.Print("Error parsing html: ", err)
		}

		essay := doc.Find(essayClass).Text()
		fmt.Println(essay)
		essaysChan <- essay
	}
	wg.Done()
}

func (f *FireflyApp) countWordsWorker(essaysChan <-chan string, wg *sync.WaitGroup) {
	for essay := range essaysChan {
		fmt.Println(essay)
	}
	wg.Done()
}

func (f *FireflyApp) fetchAndProcessEssays() {
	// open url file
	file, err := os.Open("./" + f.urlsFilename)
	if err != nil {
		fmt.Println("Error opening urls file: ", err)
		return
	}
	defer file.Close()

	// create channels for communication and synchronization
	urlsChan := make(chan string)
	essaysChan := make(chan string)

	// start a workers pool for fetching the essays
	var fetchWG sync.WaitGroup
	fetchWG.Add(numOfFetchWorkers)
	for i := 0; i < numOfFetchWorkers; i++ {
		go f.fetchEssayWorker(urlsChan, essaysChan, &fetchWG)
	}

	// start a workers pool for counting the words
	var countWG sync.WaitGroup
	countWG.Add(numOfCountWorkers)
	for i := 0; i < numOfCountWorkers; i++ {
		go f.countWordsWorker(essaysChan, &countWG)
	}

	// read urls from file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urlsChan <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	close(urlsChan)
	fetchWG.Wait()
	close(essaysChan)
	countWG.Wait()

	fmt.Println("finished")
}

// func (f *FireflyApp) fetchAndProcessEssays2() {
// file, err := os.Open("./" + f.urlsFilename)
// if err != nil {
// 	fmt.Println("Error opening urls file: ", err)
// 	return
// }
// defer file.Close()

// scanner := bufio.NewScanner(file)
// urlsChan := make(chan string)

// // insert urls to chan
// go func() {
// 	for scanner.Scan() {
// 		urlsChan <- scanner.Text()
// 	}
// 	close(urlsChan)
// }()

// wg := sync.WaitGroup{}
// wg.Add(2)
// for i := 0; i < 2; i++ {
// 	go func() {
// 		for url := range urlsChan {
// 			resp, err := f.fetcher.GetWithRetry(url, 10)
// 			if err != nil {
// 				fmt.Println("error getting url: ", err)
// 			}

// 			body, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				fmt.Print("error reading resp body: ", err)
// 			}

// 			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
// 			if err != nil {
// 				fmt.Print("Error parsing html: ", err)
// 			}

// 			essay := doc.Find(essayClass).Text()
// 			fmt.Println(essay)
// 		}
// 		wg.Done()
// 	}()
// }
// wg.Wait()
// fmt.Println("finished")
// }

func (f *FireflyApp) runApp() {
	f.createDict()
	f.fetchAndProcessEssays()
}

func main() {
	limiter := rate.NewLimiter(2, 1)

	// initialize app. I like working with interfaces to maintain a testable code
	app := FireflyApp{
		wordsFilename: "words-mini.txt",
		urlsFilename:  "endg-urls-mini",
		fetcher:       fetcher.NewFetcher(limiter),
	}
	app.runApp()
}
