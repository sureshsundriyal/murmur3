// Copyright (c) 2013, Suresh Sundriyal. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
// This is a progressive processing version of the MurmurHash3 family of hash
// functions by Austin Appleby and is a translation of the public domain code
// that can be found here:
// https://code.google.com/p/smhasher/source/browse/trunk/MurmurHash3.cpp?r=150

package murmur3

import "hash"

//Hash128 interface for 128-bit hash functions.
type Hash128 interface {
	hash.Hash
	Sum128() (uint64, uint64)
	//SetSeed sets the seed after the hash has been Reset.
	SetSeed(seed uint32) error
}

type Hash32 interface {
	hash.Hash
	Sum32() uint32
	//SetSeed sets the seed after the hash has been Reset.
	SetSeed(seed uint32) error
}
