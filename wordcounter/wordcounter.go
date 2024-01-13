package wordcounter

import (
	"container/heap"
	"errors"
	"sync"
)

type WordCounter interface {
	Increase(word string) error
	GetTopK(k int) map[string]int
}

type WordCounterImp struct {
	smap  sync.Map
	mHeap *minHeap
}

type wordCount struct {
	word  string
	count int
}

type minHeap []wordCount

func NewWordCounter() WordCounter {
	return &WordCounterImp{
		mHeap: &minHeap{},
	}
}

func (wc *WordCounterImp) Increase(word string) error {
	currentCount, ok := wc.smap.Load(word)
	if !ok {
		currentCount = 0
	}

	currentCountInt, ok := currentCount.(int)
	if !ok {
		return errors.New("word count is not a number")
	}

	wc.smap.Store(word, currentCountInt+1)
	return nil
}

// implement the heap interface
func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h minHeap) Swap(i, j int) {
	if i < 0 || i >= len(h) || j < 0 || j >= len(h) {
		return // prevent out of range access
	}
	h[i], h[j] = h[j], h[i]
}

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(wordCount))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	if n == 0 {
		return nil
	}
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (wc *WordCounterImp) GetTopK(k int) map[string]int {
	// GetTop returns a map of the K most popular words and
	// how many times they occured. We extract the top K words
	// using a designated internal minHeap

	heap.Init(wc.mHeap)

	// insert words to heap
	wc.smap.Range(func(key, value any) bool {
		heap.Push(wc.mHeap, wordCount{
			word:  key.(string),
			count: value.(int),
		})

		// maintain only x items in the heap
		if wc.mHeap.Len() > k {
			heap.Pop(wc.mHeap)
		}
		return true
	})

	topWords := make(map[string]int)
	numOfWordsToReturn := k
	if wc.mHeap.Len() < k {
		numOfWordsToReturn = wc.mHeap.Len()
	}
	for i := 0; i < numOfWordsToReturn; i++ {
		wordCountItem := heap.Pop(wc.mHeap).(wordCount)
		topWords[wordCountItem.word] = wordCountItem.count
	}

	return topWords
}
