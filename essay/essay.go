package essay

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Essay struct {
	Url    string
	client http.Client
}

const rateLimitStatusCode = 999

func (e *Essay) FetchOld() {
	// Fetch fetches an essay. If the rate limit is exceeded, we retry

	// if err := e.RateLimiter.Wait(context.Background()); err != nil {
	// 	fmt.Println("Could not wait:", err)
	// 	return
	// }

	waitTime := 1
	var res *http.Response
	var err error
	for {
		res, err = http.Get(e.Url)
		if res.StatusCode == rateLimitStatusCode {
			// fmt.Println("failed, waiting ", waitTime, " seconds")
			time.Sleep(time.Second * time.Duration(waitTime))
			waitTime *= 2
			continue
		}
		if err != nil || res.StatusCode != http.StatusOK { //?
			fmt.Println("Error fetching URL:", res.Status)
			return
		}
		break
	}
	if waitTime > 1 {
		fmt.Println(waitTime)
	}
	defer res.Body.Close()
	_, err = io.ReadAll(res.Body)
	if err != nil { //?
		fmt.Println("Error reading res body:", err)
		return
	}

	// fmt.Println(string(body)[:10])

	// res, err := http.Get(e.Url)
	// if err != nil || res.StatusCode != http.StatusOK { //?
	// 	fmt.Println("Error fetching URL:", res.Status)
	// 	return
	// }
	// defer res.Body.Close()
	// body, err := io.ReadAll(res.Body)
	// if err != nil { //?
	// 	fmt.Println("Error reading res body:", err)
	// 	return
	// }

	// fmt.Println(string(body)[:10])

}

func (e *Essay) Fetch() {
	// Fetch fetches an essay. If the rate limit is exceeded, we retry

}
