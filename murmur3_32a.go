// Copyright (c) 2013, Suresh Sundriyal. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
// This is a progressive processing version of the MurmurHash3 family of hash
// functions by Austin Appleby and is a translation of the public domain code
// that can be found here:
// https://code.google.com/p/smhasher/source/browse/trunk/MurmurHash3.cpp?r=150

package murmur3

import "errors"

const (
	//Constants for x86 32-bit hash function.
	c1_32_32 = 0xcc9e2d51
	c2_32_32 = 0x1b873593
)

//sum32_32 struct contains variables used in x86 32-bit hash calculations.
type sum32_32 struct {
	h1     uint32
	k1     uint32
	length uint32
	offset uint8
}

// New32a returns a Murmur3 32-bit hash.Hash opmtimized for 32-bit architecture.
func New32a(seed uint32) Hash32 {
	return &sum32_32{seed, 0, 0, 0}
}

// Reset resets the hash to one with zero bytes written.
func (s *sum32_32) Reset() {
	s.h1, s.k1, s.length, s.offset = 0, 0, 0, 0
}

func (s *sum32_32) SetSeed(seed uint32) error {
	if s.h1 != 0 {
		return errors.New("hash needs to be reset")
	} else {
		s.h1 = seed
		return nil
	}
}

func (s *sum32_32) Write(data []byte) (int, error) {
	length := len(data)
	if length == 0 {
		return 0, nil
	}
	s.length += uint32(length)

	for _, c := range data {
		// TODO: Might want to check this for endianness for consistency
		// across systems.
		if s.offset < 4 {
			s.k1 |= uint32(uint32(c) << uint32(s.offset*8))
		}
		s.offset++

		if s.offset == 4 {
			s.k1 *= c1_32_32
			s.k1 = (s.k1 << 15) | (s.k1 >> (32 - 15))
			s.k1 *= c2_32_32

			s.h1 ^= s.k1
			s.h1 = (s.h1 << 13) | (s.h1 >> (32 - 13))
			s.h1 = s.h1*5 + 0xe6546b64

			s.k1 = 0
			s.offset = 0
		}
	}
	return length, nil
}

func (s *sum32_32) Sum32() uint32 {
	var h1 = s.h1
	var k1 = s.k1

	//tail
	if k1 != 0 {
		k1 *= c1_32_32
		k1 = (k1 << 16) | (k1 >> (32 - 16))
		k1 *= c2_32_32
		h1 ^= k1
	}

	//finalization
	h1 ^= s.length

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}

func (s *sum32_32) Sum(in []byte) []byte {
	h1 := s.Sum32()
	return append(in, byte(h1>>24), byte(h1>>16), byte(h1>>8), byte(h1))
}

func (s *sum32_32) BlockSize() int { return 4 }

func (s *sum32_32) Size() int { return 4 }
