package dictionary

import (
	"sync"

	"github.com/spaolacci/murmur3"
)

type Dictionary struct {
	sync.RWMutex
	tokens map[uint64]int
	values [][]byte
}

func NewDictionary() *Dictionary {
	tokens := make(map[uint64]int)
	values := make([][]byte, 1)
	return &Dictionary{tokens: tokens, values: values}
}

func (dict *Dictionary) token(hash uint64) (int, bool) {
	dict.RLock()
	token, ok := dict.tokens[hash]
	dict.RUnlock()
	return token, ok
}

func (dict *Dictionary) AddToken(value []byte) (int, bool) {
	hash := murmur3.Sum64(value)
	if token, ok := dict.token(hash); ok {
		return token, true
	}
	dict.Lock()
	if token, ok := dict.tokens[hash]; ok {
		dict.Unlock()
		return token, true
	}
	token := len(dict.values)
	temp := make([]byte, len(value))
	copy(temp, value)
	dict.values = append(dict.values, temp)
	dict.tokens[hash] = token
	dict.Unlock()
	return token, false
}

func (dict *Dictionary) Token(value []byte) (int, bool) {
	hash := murmur3.Sum64(value)
	return dict.token(hash)
}

func (dict *Dictionary) Value(token int) ([]byte, bool) {
	if token > 0 && token < len(dict.values) {
		return dict.values[token], true
	}
	return nil, false
}

func (dict *Dictionary) Len() int {
	return len(dict.values)
}
