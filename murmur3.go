// Copyright (c) 2013, Suresh Sundriyal. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
// This is a progressive processing version of the MurmurHash3 family of hash
// functions by Austin Appleby and is a translation of the public domain code
// that can be found here:
// https://code.google.com/p/smhasher/source/browse/trunk/MurmurHash3.cpp?r=150

package murmur3

import "hash"

//Extend hash.Hash to accomodate for setting the salt in Reset()
type HashM3 interface {
	hash.Hash
	//ResetAndSetSeed resets the hash and sets the seed.
	ResetAndSetSeed(seed uint32)
}

type Hash32 interface {
	HashM3
	Sum32() uint32
}

//Hash128 interface for 128-bit hash functions.
type Hash128 interface {
	HashM3
	Sum128() (uint64, uint64)
}
