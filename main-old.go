package main

// import (
// 	"bufio"
// 	"firefly/essay"
// 	"fmt"
// 	"os"
// 	"regexp"
// 	"strings"
// 	"sync"
// 	"time"
// )

// type Firefly struct {
// 	urlsFilename  string
// 	wordsFilename string
// 	urlsChan      chan string
// 	wordsChan     chan string
// 	createdDictChan chan bool
// 	numOfWorkers  int
// 	chunkSize     int
// 	wordsDict     sync.Map
// }

// func (f *Firefly) loadUrls() {
// 	file, err := os.Open("./" + f.urlsFilename)
// 	if err != nil {
// 		fmt.Println("Error opening file: ", err)
// 		return
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		f.urlsChan <- scanner.Text()
// 	}

// 	close(f.urlsChan)
// }

// func (f *Firefly) fetchEssay(i int, wg *sync.WaitGroup) {
// 	for url := range f.urlsChan {
// 		// fmt.Println("worker ", i, " got url: ", url)
// 		e := essay.Essay{
// 			Url: url,
// 		}
// 		e.Fetch()
// 	}
// 	wg.Done()
// }

// func isAlphabetic2(s string) bool {
// 	return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(s)
// }

// func (f *Firefly) processWordsWorker(wg *sync.WaitGroup) {
// 	for word := range f.wordsChan {
// 		if len(word) < 3 || !isAlphabetic(word) {
// 			continue
// 		}
// 		f.wordsDict.Store(strings.ToLower(word), true)
// 	}
// 	wg.Done()
// }

// func (f *Firefly) createWordsDict() {
// 	file, err := os.Open("./" + f.wordsFilename)
// 	if err != nil {
// 		fmt.Println("Error opening file: ", err)
// 		return
// 	}
// 	defer file.Close()

// 	// initialize word dict
// 	f.wordsDict = sync.Map{}
// 	f.wordsChan = make(chan string)
// 	scanner := bufio.NewScanner(file)

// 	go func() {
// 		for scanner.Scan() {
// 			f.wordsChan <- scanner.Text()
// 		}

// 		close(f.wordsChan)
// 	}()

// 	wg := sync.WaitGroup{}
// 	wg.Add(f.numOfWorkers)
// 	for i := 0; i < f.numOfWorkers; i++ {
// 		go f.processWordsWorker(&wg)
// 	}
// 	wg.Wait()
// }

// func (f *Firefly) RunApp() {
// 	fmt.Println("running")
// 	start := time.Now()

// 	// Create a dict of valid words
// 	f.createWordsDict()

// 	f.urlsChan = make(chan string)
// 	go f.loadUrls()

// 	// limiter := rate.NewLimiter(100, 3)
// 	wg := sync.WaitGroup{}
// 	for i := 0; i < f.numOfWorkers; i++ {
// 		wg.Add(1)
// 		go f.fetchEssay(i, &wg)
// 	}
// 	wg.Wait()
// 	elapsed := time.Since(start)
// 	fmt.Printf("Processes took %s", elapsed)
// }

// func mainOld() {
// 	firefly := Firefly{
// 		urlsFilename:  "endg-urls-mini",
// 		wordsFilename: "words.txt",
// 		numOfWorkers:  15,
// 		chunkSize:     1000,
// 	}

// 	firefly.RunApp()
// }
