package dictionary

import (
	"sync"

	"github.com/spaolacci/murmur3"
)

type Dictionary struct {
	sync.RWMutex
	hashes map[uint64]int
	values [][]byte
}

func NewDictionary() *Dictionary {
	hashes := make(map[uint64]int)
	values := make([][]byte, 1)
	return &Dictionary{hashes: hashes, values: values}
}

func (dict *Dictionary) Identify(val []byte) (int, error) {
	hash := murmur3.Sum64(val)

	dict.RLock()
	if id, ok := dict.hashes[hash]; ok {
		dict.RUnlock()
		return id, nil
	}
	dict.RUnlock()

	dict.Lock()
	if id, ok := dict.hashes[hash]; ok {
		dict.Unlock()
		return id, nil
	}
	last := len(dict.values)
	temp := make([]byte, len(val))
	copy(temp, val)
	temp = unquote(temp)
	dict.values = append(dict.values, temp)
	dict.hashes[hash] = last
	if len(temp) < len(val) {
		hash = murmur3.Sum64(temp)
		dict.hashes[hash] = last
	}
	dict.Unlock()
	return last, nil
}

func (dict *Dictionary) Len() int {
	return len(dict.values)
}

func (dict *Dictionary) Value(id int) ([]byte, bool) {
	if id > 0 && id < len(dict.values) {
		return dict.values[id], true
	}
	return nil, false
}
