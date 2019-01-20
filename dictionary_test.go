package dictionary

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDictionary(t *testing.T) {
	dict := NewDictionary()
	for i := 1; i <= 50; i++ {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(i))
		id, err := dict.Identify(buf)
		if err != nil {
			t.Fail()
			return
		}
		val, ok := dict.Value(id)
		if !ok || !bytes.Equal(buf, val) {
			t.Fail()
			return
		}
	}
}

func TestDictionarySize(t *testing.T) {
	dict := NewDictionary()
	dict.Identify([]byte{0x10})
	dict.Identify([]byte{0x20})
	dict.Identify([]byte{0x30})
	_, ok := dict.Value(100)
	if ok {
		t.Fail()
		return
	}
}

func TestIdentify(t *testing.T) {
	dict := NewDictionary()
	if a, err := dict.Identify([]byte{1, 2, 3}); err != nil {
		t.Fail()
	} else {
		t.Log(a)
	}
	if a, err := dict.Identify([]byte{1, 2, 3}); err != nil {
		t.Fail()
	} else {
		t.Log(a)
	}
	if a, err := dict.Identify([]byte{1, 2, 3}); err != nil {
		t.Fail()
	} else {
		t.Log(a)
	}
}
