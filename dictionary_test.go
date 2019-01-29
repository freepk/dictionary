package dictionary

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDictionary(t *testing.T) {
	dict := NewDictionary(1)
	for i := 1; i <= 50; i++ {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(i))
		key, ok := dict.AddKey(buf)
		if ok {
			t.Fail()
			return
		}
		val, ok := dict.Val(key)
		if !ok || !bytes.Equal(buf, val) {
			t.Fail()
			return
		}
	}
}
