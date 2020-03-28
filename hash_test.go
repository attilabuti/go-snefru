package snefru

import (
	"encoding/hex"
	"testing"
)

var vec = []string{"hello", "7c5f22b1a92d9470efea37ec6ed00b2357a4ce3c41aa6e28e3b84057465dbb56"}

func TestNewSnefru256(t *testing.T) {
	h := NewSnefru256(8)
	h.Write([]byte(vec[0]))
	s := hex.EncodeToString(h.Sum(nil))
	if s != vec[1] {
		t.Log(s)
		t.Fail()
	}
}
