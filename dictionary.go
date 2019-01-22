package dictionary

import (
	"sync"
)

type Dictionary struct {
	sync.RWMutex
	tokens map[string]int
	values [][]byte
}

func NewDictionary() *Dictionary {
	tokens := make(map[string]int)
	values := make([][]byte, 1)
	return &Dictionary{tokens: tokens, values: values}
}

func (dict *Dictionary) AddToken(value []byte) (int, bool) {
	dict.Lock()
	if token, ok := dict.tokens[string(value)]; ok {
		dict.Unlock()
		return token, true
	}
	token := len(dict.values)
	temp := make([]byte, len(value))
	copy(temp, value)
	dict.values = append(dict.values, temp)
	dict.tokens[string(value)] = token
	dict.Unlock()
	return token, false
}

func (dict *Dictionary) Token(value []byte) (int, bool) {
	dict.RLock()
	token, ok := dict.tokens[string(value)]
	dict.RUnlock()
	return token, ok
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
