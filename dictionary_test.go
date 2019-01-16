package dictionary

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDictionary(t *testing.T) {
	dict, err := NewDictionary(256)
	if err != nil {
		t.Fail()
		return
	}

	for i := 1; i <= 50; i++ {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(i))
		id, err := dict.Identify(buf)
		if err != nil {
			t.Fail()
			return
		}
		val, err := dict.Value(id)
		if err != nil || !bytes.Equal(buf, val) {
			t.Fail()
			return
		}
	}
}

func TestDictionarySize(t *testing.T) {
	dict, err := NewDictionary(4)
	if err != nil {
		t.Fail()
		return
	}
	dict.Identify([]byte{0x10})
	dict.Identify([]byte{0x20})
	dict.Identify([]byte{0x30})
	_, err = dict.Identify([]byte{0x40})
	if err != SizeOverflowError {
		t.Fail()
		return
	}
	_, err = dict.Value(100)
	if err != NotExistsError {
		t.Fail()
		return
	}
}
