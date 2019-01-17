package dictionary

import (
	"errors"
	"sync/atomic"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
)

var (
	SizeOverflowError = errors.New("Size overflow")
	NotExistsError    = errors.New("Not exists")
)

type Dictionary struct {
	last   uint64
	hits   []uint64
	hashes *hashtab.HashTab
	values [][]byte
}

func NewDictionary(num int) *Dictionary {
	hashes := hashtab.NewHashTab(num)
	hits := make([]uint64, num)
	values := make([][]byte, num)
	return &Dictionary{last: 1, hits: hits, hashes: hashes, values: values}
}

func (dict *Dictionary) Identify(val []byte) (uint64, error) {
	hash := murmur3.Sum64(val)
	if id, ok := dict.hashes.Get(hash); ok {
		atomic.AddUint64(&dict.hits[id], 1)
		return id, nil
	}
	last := atomic.LoadUint64(&dict.last)
	if int(last) > len(dict.values) - 1 {
		return 0, SizeOverflowError
	}
	id, ok := dict.hashes.GetOrSet(hash, last)
	if ok {
		atomic.AddUint64(&dict.hits[id], 1)
		return id, nil
	}
	atomic.AddUint64(&dict.last, 1)
	tmp := make([]byte, len(val))
	copy(tmp, val)
	dict.values[id] = tmp
	return id, nil
}

func (dict *Dictionary) Value(id uint64) ([]byte, error) {
	if id > 0 && id < dict.last {
		return dict.values[id], nil
	}
	return nil, NotExistsError
}

func (dict *Dictionary) Hits() []uint64 {
	return dict.hits
}

