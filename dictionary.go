package dictionary

import (
	"sync"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
)

type Dictionary struct {
	sync.RWMutex
	lastKey int
	keys    *hashtab.HashTab
	values  [][]byte
	hits    []int
}

func NewDictionary(power uint8) *Dictionary {
	keys := hashtab.NewHashTab(power)
	if keys == nil {
		return nil
	}
	values := make([][]byte, keys.Size()+1)
	hits := make([]int, keys.Size()+1)
	return &Dictionary{lastKey: 1, keys: keys, values: values, hits: hits}
}

func (d *Dictionary) GetKey(val []byte) (int, bool) {
	hash := murmur3.Sum64(val)
	// Fast path
	d.RLock()
	key, ok := d.keys.Get(hash)
	if ok {
		d.RUnlock()
		return int(key), true
	}
	d.RUnlock()
	// Slow path
	d.Lock()
	key, ok = d.keys.Get(hash)
	if ok {
		d.Unlock()
		return int(key), true
	}
	key = uint64(d.lastKey)
	d.keys.Set(hash, key)
	d.lastKey++
	d.hits[key]++
	d.Unlock()
	// Copy value after key assigment
	temp := make([]byte, len(val))
	copy(temp, val)
	d.values[key] = temp

	return int(key), false
}

func (d *Dictionary) GetValue(key int) ([]byte, bool) {
	if key > d.lastKey {
		return nil, false
	}
	return d.values[key], true
}
