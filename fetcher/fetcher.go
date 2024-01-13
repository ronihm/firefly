package fetcher

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"golang.org/x/time/rate"
)

type Fetcher interface {
	Get(url string) (*http.Response, error)
	GetWithRetry(url string, maxRetries int) (*http.Response, error)
}

type FetcherImp struct {
	RateLimiter *rate.Limiter
	client      *http.Client
}

func NewFetcher(rateLimiter *rate.Limiter) Fetcher {
	f := FetcherImp{
		RateLimiter: rateLimiter,
		client:      &http.Client{},
	}

	return &f
}

const rateLimitStatusCode = 999

func (f *FetcherImp) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error creating GET request", err)
		return nil, err
	}

	// set random user agent
	req.Header.Set("User-Agent", browser.Random())

	if err := f.RateLimiter.Wait(context.Background()); err != nil {
		fmt.Println("could not wait:", err)
		return nil, err
	}

	return f.client.Do(req)

}

func (f *FetcherImp) GetWithRetry(url string, maxRetries int) (*http.Response, error) {
	// Fetch fetches a url. If the rate limit is exceeded, we retry
	var err error

	for i := 0; i < maxRetries; i++ {
		resp, err := f.Get(url)
		if resp.StatusCode >= 400 && resp.StatusCode != rateLimitStatusCode {
			// if legitimate error e.g. 404, return
			return resp, err
		} else if err == nil {
			return resp, nil
		}

		wait := math.Pow(2, float64(i)) // Exponential backoff
		fmt.Printf("Attempt %d failed, retrying in %v seconds...\n", i+1, wait)
		time.Sleep(time.Duration(wait) * time.Second)
	}

	return nil, err
}
