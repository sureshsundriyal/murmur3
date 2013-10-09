package murmur3

import (
	"testing"
)

func TestAll(t *testing.T) {
	s := []byte("hello")

	h128 := New64_128()

	h128.Write(s)

	h1, h2 := h128.Sum128()

	if h1 != uint64(14688674573012802306) && h2 != uint64(6565844092913065241) {
		t.Error("Something is wrong here.")
	}
}
