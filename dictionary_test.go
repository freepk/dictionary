package dictionary

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDictionary(t *testing.T) {
	dict := NewDictionary(8)
	if dict == nil {
		t.Fail()
		return
	}
	for i := 1; i <= 50; i++ {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(i))
		key, err := dict.GetKey(buf)
		if err != nil {
			t.Fail()
			return
		}
		val, err := dict.GetValue(key)
		if err != nil || !bytes.Equal(buf, val) {
			t.Fail()
			return
		}
	}
}

func TestDictionarySize(t *testing.T) {
	dict := NewDictionary(2)
	if dict == nil {
		t.Fail()
		return
	}
	dict.GetKey([]byte{0x10})
	dict.GetKey([]byte{0x20})
	dict.GetKey([]byte{0x30})
	dict.GetKey([]byte{0x40})
}
