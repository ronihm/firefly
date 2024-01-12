package safetrie

import (
	"sync"

	"github.com/dghubble/trie"
)

type SafeTrie struct {
	trie *trie.PathTrie
	mut  sync.Mutex
}

func NewSafeTrie() *SafeTrie {
	return &SafeTrie{
		trie: trie.NewPathTrie(),
	}
}

func (t *SafeTrie) Put(key string, value bool) bool {
	t.mut.Lock()
	defer t.mut.Unlock()

	return t.trie.Put(key, value)
}

func (t *SafeTrie) IsInTrie(key string) bool {
	t.mut.Lock()
	defer t.mut.Unlock()

	if t.trie.Get(key) == nil {
		return false
	}

	return true
}
