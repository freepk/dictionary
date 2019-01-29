package dictionary

import (
	"sync"

	"github.com/spaolacci/murmur3"
)

type Dictionary struct {
	sync.RWMutex
	keys map[uint64]int
	vals [][]byte
}

func NewDictionary(reserved int) *Dictionary {
	keys := make(map[uint64]int)
	vals := make([][]byte, reserved)
	return &Dictionary{keys: keys, vals: vals}
}

func (dict *Dictionary) AddKey(val []byte) (int, bool) {
	hash := murmur3.Sum64(val)
	dict.RLock()
	if key, ok := dict.keys[hash]; ok {
		dict.RUnlock()
		return key, true
	}
	dict.RUnlock()
	dict.Lock()
	if key, ok := dict.keys[hash]; ok {
		dict.Unlock()
		return key, true
	}
	key := len(dict.vals)
	dict.keys[hash] = key
	dict.vals = append(dict.vals, append([]byte{}, val...))
	dict.Unlock()
	return key, false
}

func (dict *Dictionary) Key(val []byte) (int, bool) {
	hash := murmur3.Sum64(val)
	dict.RLock()
	key, ok := dict.keys[hash]
	dict.RUnlock()
	return key, ok
}

func (dict *Dictionary) Val(key int) ([]byte, bool) {
	if key > 0 && key < len(dict.vals) {
		return dict.vals[key], true
	}
	return nil, false
}

func (dict *Dictionary) Len() int {
	return len(dict.vals)
}
