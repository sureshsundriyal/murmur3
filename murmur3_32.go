// Copyright (c) 2013, Suresh Sundriyal. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
// This is a progressive processing version of the MurmurHash3 family of hash
// functions by Austin Appleby and is a translation of the public domain code
// that can be found here:
// https://code.google.com/p/smhasher/source/browse/trunk/MurmurHash3.cpp?r=150

package murmur3

const (
	// Constants for x86 128-bit hash function.
	c1_32_128 = 0x239b961b
	c2_32_128 = 0xab0e9789
	c3_32_128 = 0x38b34ae5
	c4_32_128 = 0xa1e38b93
)

// sum32_128 struct contains variables used in x86 128-bit hash calculations.
type sum32_128 struct {
	h1     uint32
	h2     uint32
	h3     uint32
	h4     uint32
	k1     uint32
	k2     uint32
	k3     uint32
	k4     uint32
	length uint32
	offset uint8
}

// New32 returns a Murmur3 128-bit hash.Hash optimized for 32-bit architecture.
func New32(seed uint32) Hash128 {
	return &sum32_128{seed, seed, seed, seed, seed, 0, 0, 0, 0, 0}
}

// Reset resets the hash to one with zero bytes written.
func (s *sum32_128) Reset() {
	s.h1, s.h2, s.h3, s.h4 = 0, 0, 0, 0
	s.k1, s.k2, s.k3, s.k4 = 0, 0, 0, 0
	s.length, s.offset = 0, 0
}

func (s *sum32_128) ResetAndSetSeed(seed uint32) {
	s.Reset()
	s.h1, s.h2, s.h3, s.h4 = seed, seed, seed, seed
}

func (s *sum32_128) Sum128() (uint64, uint64) {
	var h1, h2, h3, h4 = s.h1, s.h2, s.h3, s.h4
	var k1, k2, k3, k4 uint32 = s.k1, s.k2, s.k3, s.k4

	//tail
	switch {
	case s.offset > 12:
		k4 *= c4_32_128
		k4 = (k4 << 18) | (k4 >> (32 - 18))
		k4 *= c1_32_128
		h4 ^= k4
		fallthrough

	case s.offset > 8:
		k3 *= c3_32_128
		k3 = (k3 << 17) | (k3 >> (32 - 17))
		k3 *= c4_32_128
		h3 ^= k3
		fallthrough

	case s.offset > 4:
		k2 *= c2_32_128
		k2 = (k2 << 16) | (k2 >> (32 - 16))
		k2 *= c3_32_128
		h2 ^= k2
		fallthrough

	case s.offset > 0:
		k1 *= c1_32_128
		k1 = (k1 << 15) | (k1 >> (32 - 15))
		k1 *= c2_32_128
		h1 ^= k1
	}

	// finalization
	h1 ^= s.length
	h2 ^= s.length
	h3 ^= s.length
	h4 ^= s.length

	h1 += h2
	h1 += h3
	h1 += h4

	h2 += h1
	h3 += h1
	h4 += h1

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	h2 ^= h2 >> 16
	h2 *= 0x85ebca6b
	h2 ^= h2 >> 13
	h2 *= 0xc2b2ae35
	h2 ^= h2 >> 16

	h3 ^= h3 >> 16
	h3 *= 0x85ebca6b
	h3 ^= h3 >> 13
	h3 *= 0xc2b2ae35
	h3 ^= h3 >> 16

	h4 ^= h4 >> 16
	h4 *= 0x85ebca6b
	h4 ^= h4 >> 13
	h4 *= 0xc2b2ae35
	h4 ^= h4 >> 16

	h1 += h2
	h1 += h3
	h1 += h4

	h2 += h1
	h3 += h1
	h4 += h1

	return uint64((uint64(h2) << 32) | uint64(h1)),
		uint64((uint64(h4) << 32) | uint64(h3))
}

func (s *sum32_128) Sum(in []byte) []byte {
	h1, h2 := s.Sum128()
	return append(in, byte(h1>>56), byte(h1>>48), byte(h1>>40), byte(h1>>32),
		byte(h1>>24), byte(h1>>16), byte(h1>>8), byte(h1), byte(h2>>56),
		byte(h2>>48), byte(h2>>32), byte(h2>>24), byte(h2>>16),
		byte(h2>>8), byte(h2))
}

func (s *sum32_128) Write(data []byte) (int, error) {
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
		} else if s.offset >= 4 && s.offset < 8 {
			s.k2 |= uint32(uint32(c) << uint32((s.offset%4)*8))
		} else if s.offset >= 8 && s.offset < 12 {
			s.k3 |= uint32(uint32(c) << uint32((s.offset%4)*8))
		} else if s.offset >= 12 && s.offset < 16 {
			s.k4 |= uint32(uint32(c) << uint32((s.offset%4)*8))
		}
		s.offset++

		if s.offset == 16 {
			s.k1 *= c1_32_128
			s.k1 = (s.k1 << 15) | (s.k1 >> (32 - 15))
			s.k1 *= c2_32_128
			s.h1 ^= s.k1

			s.h1 = (s.h1 << 19) | (s.h1 >> (32 - 19))
			s.h1 += s.h2
			s.h1 = s.h1*5 + 0x561ccd1b

			s.k2 *= c2_32_128
			s.k2 = (s.k2 << 16) | (s.k2 >> (32 - 16))
			s.k2 *= c3_32_128
			s.h2 ^= s.k2

			s.h2 = (s.h2 << 17) | (s.h2 >> (32 - 17))
			s.h2 += s.h3
			s.h2 = s.h2*5 + 0x0bcaa747

			s.k3 *= c3_32_128
			s.k3 = (s.k3 << 17) | (s.k3 >> (32 - 17))
			s.k3 *= c4_32_128
			s.h3 ^= s.k3

			s.h3 = (s.h3 << 15) | (s.h3 >> (32 - 15))
			s.h3 += s.h4
			s.h3 = s.h3*5 + 0x96cd1c35

			s.k4 *= c4_32_128
			s.k4 = (s.k4 << 18) | (s.k4 >> (32 - 18))
			s.k4 *= c1_32_128
			s.h4 ^= s.k4

			s.h4 = (s.h4 << 13) | (s.h4 >> (32 - 13))
			s.h4 += s.h1
			s.h4 = s.h4*5 + 0x32ac3b17

			s.k1, s.k2, s.k3, s.k4 = 0, 0, 0, 0
			s.offset = 0
		}
	}
	return length, nil
}

func (s *sum32_128) BlockSize() int { return 16 }

func (s *sum32_128) Size() int { return 16 }
