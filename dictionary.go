package dictionary

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
)

var (
	SizeOverflowError = errors.New("Size overflow")
	NotExistsError    = errors.New("Not exists")
)

type Dictionary struct {
	sync.Mutex
	size   uint64
	last   uint64
	hits   []uint64
	hashes *hashtab.HashTab
	values [][]byte
}

func NewDictionary(size uint32) *Dictionary {
	hashes := hashtab.NewHashTab(size)
	hits := make([]uint64, size)
	values := make([][]byte, size)
	return &Dictionary{size: uint64(size), last: 1, hits: hits, hashes: hashes, values: values}
}

func (dict *Dictionary) Identify(val []byte) (uint64, error) {
	hash := murmur3.Sum64(val)
	if id, ok := dict.hashes.Get(hash); ok {
		atomic.AddUint64(&dict.hits[id], 1)
		return id, nil
	}
	tmp := make([]byte, len(val))
	copy(tmp, val)
	dict.Lock()
	id, ok := dict.hashes.Get(hash)
	if ok {
		dict.Unlock()
		atomic.AddUint64(&dict.hits[id], 1)
		return id, nil
	}
	if dict.last >= dict.size {
		dict.Unlock()
		return 0, SizeOverflowError
	}
	id = dict.last
	dict.last++
	dict.values[id] = tmp
	dict.Unlock()
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
