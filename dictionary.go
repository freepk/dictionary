package dictionary

import (
	"errors"
	"sync"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
)

var (
	OverflowError     = errors.New("Overflow")
	KeyNotExistsError = errors.New("Key not exists")
)

type Dictionary struct {
	sync.RWMutex
	lastKey int
	keys    *hashtab.HashTab
	values  [][]byte
}

func NewDictionary(power uint8) *Dictionary {
	keys := hashtab.NewHashTab(power)
	if keys == nil {
		return nil
	}
	values := make([][]byte, keys.Size())
	return &Dictionary{lastKey: 1, keys: keys, values: values}
}

func (d *Dictionary) GetKey(val []byte) (int, error) {
	hash := murmur3.Sum64(val)
	// Fast path
	d.RLock()
	key, ok := d.keys.Get(hash)
	if ok {
		d.RUnlock()
		return int(key), nil
	}
	d.RUnlock()
	// Slow path
	d.Lock()
	if d.lastKey == int(d.keys.Size()) {
		return 0, OverflowError
	}
	key, ok = d.keys.Get(hash)
	if ok {
		d.Unlock()
		return int(key), nil
	}
	key = uint64(d.lastKey)
	d.keys.Set(hash, key)
	d.lastKey++
	d.Unlock()
	// Copy value after key assigment
	temp := make([]byte, len(val))
	copy(temp, val)
	d.values[key] = temp

	return int(key), nil
}

func (d *Dictionary) GetValue(key int) ([]byte, error) {
	if key > 0 && key < d.lastKey {
		return d.values[key], nil
	}
	return nil, KeyNotExistsError
}
