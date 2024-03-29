package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"firefly/fetcher"
	"firefly/safetrie"
	"firefly/wordcounter"

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
	wordsCounter  wordcounter.WordCounter
}

func isAlphabetic(s string) bool {
	return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(s)
}

func (f *FireflyApp) createWordsBank() {
	fmt.Println("creating words bank")
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
		essaysChan <- essay
	}
	wg.Done()
}

func (f *FireflyApp) countWordsWorker(essaysChan <-chan string, wg *sync.WaitGroup) {
	// ignore punctuation marks like brackets, quotes, etc...
	isSeparator := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && string(c) != "'"
	}

	for essay := range essaysChan {
		words := strings.FieldsFunc(essay, isSeparator)
		for _, word := range words {
			word = strings.ToLower(word)
			if f.wordsBank.IsInTrie(word) {
				f.wordsCounter.Increase(word)
			}
		}
	}
	wg.Done()
}

func (f *FireflyApp) fetchAndProcessEssays() {
	fmt.Println("fetching and processing essays")

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

	topWords := f.wordsCounter.GetTopK(10)
	prettyJson, err := json.MarshalIndent(topWords, "", "    ")
	if err != nil {
		fmt.Println("error marshaling result: ", err)
	}
	fmt.Println(string(prettyJson))
}

func (f *FireflyApp) runApp() {
	f.createWordsBank()
	f.fetchAndProcessEssays()
}

func shouldUseMini() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", "would you like to use a shorter version of this program?")

		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" || input == "yes" {
			return true, nil
		} else if input == "n" || input == "no" {
			return false, nil
		}

		fmt.Println("Please enter 'y' for yes or 'n' for no.")
	}
}

func main() {
	limiter := rate.NewLimiter(2, 1)

	urlsFile := "endg-urls"
	wordsFile := "words.txt"

	isMini, err := shouldUseMini()
	if err != nil {
		fmt.Println("error running app: ", err)
	}
	if isMini {
		urlsFile = "endg-urls-mini"
		wordsFile = "words-mini.txt"
	}

	// initialize app. I like working with interfaces to maintain a testable code
	app := FireflyApp{
		wordsFilename: wordsFile,
		urlsFilename:  urlsFile,
		fetcher:       fetcher.NewFetcher(limiter),
		wordsCounter:  wordcounter.NewWordCounter(),
	}
	app.runApp()
}
