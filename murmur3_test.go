// Copyright (c) 2013, Suresh Sundriyal. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package murmur3_test

import (
	"testing"
	. "../murmur3"
)

func TestAll(t *testing.T) {
	s := []byte("hello")

	x := []byte(`Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.`)

	h128 := New64_128(0)
	h128.Write(s)
	h1, h2 := h128.Sum128()

	if h1 != 14688674573012802306 || h2 != 6565844092913065241 {
		t.Error("x86_64: ", s, h1, h2)
	}

	h128.Reset()
	h128.SetSeed(12345)
	h128.Write(s)
	h1, h2 = h128.Sum128()

	if h1 != 17440987278262125697 || h2 != 15376406881033980724 {
		t.Error("x86_64(seed): ", s, h1, h2)
	}

	h128.Reset()
	h128.Write(x)
	h1, h2 = h128.Sum128()

	if h1 != 1706326840306453215 || h2 != 5127165288307402704 {
		t.Error("x86_64: ", x, h1, h2)
	}

	h128 = New32_128(0)
	h128.Write(s)
	h1, h2 = h128.Sum128()

	if h1 != 15821672119091348640 || h2 != 11158567162092401078 {
		t.Error("x86: ", s, h1, h2)
	}

	h128.Reset()
	h128.Write(x)
	h1, h2 = h128.Sum128()

	if h1 != 223949659430422294 || h2 != 10022274208940483369 {
		t.Error("x86: ", x, h1, h2)
	}

	h128.Reset()
	h128.SetSeed(12345)
	h128.Write(s)
	h1, h2 = h128.Sum128()

	if h1 != 7882561715466346695 || h2 != 11883514271246235972 {
		t.Error("x86 (seed): ", s, h1, h2)
	}

}
