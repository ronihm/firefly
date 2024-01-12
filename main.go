package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"firefly/safetrie"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"
)

type FireflyApp struct {
	wordsFilename string
	urlsFilename  string
	dict          map[string]bool
	wordsBank     *safetrie.SafeTrie
}

func isAlphabetic(s string) bool {
	return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(s)
}

func (f *FireflyApp) createDict() {
	f.dict = make(map[string]bool)
	f.wordsBank = safetrie.NewSafeTrie()
	file, err := os.Open("./" + f.wordsFilename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) > 2 && isAlphabetic(word) {
			f.dict[word] = true
		}
	}
}

func (f *FireflyApp) fetchEssays() {
	file, err := os.Open("./" + f.urlsFilename)
	if err != nil {
		fmt.Println("Error opening urls file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	urlsChan := make(chan string)

	// insert urls to chan
	go func() {
		for scanner.Scan() {
			urlsChan <- scanner.Text()
		}
		close(urlsChan)
	}()

	client := &http.Client{}
	wg := sync.WaitGroup{}
	limiter := rate.NewLimiter(2, 1)
	var counter int32

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			for url := range urlsChan {
				if err := limiter.Wait(context.Background()); err != nil {
					fmt.Println("Could not wait:", err)
					return
				}
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					fmt.Println("error ", err)
				}
				req.Header.Set("User-Agent", browser.Random())
				resp, err := client.Do(req)
				if err != nil || resp.StatusCode == 999 {
					fmt.Println(err, " ", resp.Status)
				} else {
					atomic.AddInt32(&counter, 1)
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						fmt.Print("Error reading body")
					}
					doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
					if err != nil {
						fmt.Print("Error parsing html")
					}
					content := doc.Find(".caas-body").Text()
					fmt.Println(content)
				}
				fmt.Println(atomic.LoadInt32(&counter))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("finished")
}

func (f *FireflyApp) runApp() {
	f.createDict()
	f.fetchEssays()
}

func main() {
	app := FireflyApp{
		wordsFilename: "words-mini.txt",
		urlsFilename:  "endg-urls-mini",
	}
	app.runApp()
}
