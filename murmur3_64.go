// Copyright (c) 2013, Suresh Sundriyal. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
// This is a progressive processing version of the MurmurHash3 family of hash
// functions by Austin Appleby and is a translation of the public domain code
// that can be found here:
// https://code.google.com/p/smhasher/source/browse/trunk/MurmurHash3.cpp?r=150

package murmur3

const (
	//Constants for x86_64 128-bit hash function.
	c1_64_128 = 0x87c37b91114253d5
	c2_64_128 = 0x4cf5ad432745937f
)

//sum64_128 struct contains variables used in x86_64 128-bit hash calculations.
type sum64_128 struct {
	h1     uint64
	h2     uint64
	k1     uint64
	k2     uint64
	length uint64
	offset uint8
}

// New64 returns a Murmur3 128-bit hash.Hash optimized for 64-bit architecture.
func New64(seed uint32) Hash128 {
	seed64 := uint64(seed)
	return &sum64_128{seed64, seed64, 0, 0, 0, 0}
}

//Reset resets the hash to one with zero bytes written.
func (s *sum64_128) Reset() {
	s.h1, s.h2, s.k1, s.k2, s.length, s.offset = 0, 0, 0, 0, 0, 0
}

func (s *sum64_128) ResetAndSetSeed(seed uint32) {
	s.Reset()
	s.h1, s.h2 = uint64(seed), uint64(seed)
}

func (s *sum64_128) Write(data []byte) (int, error) {
	length := len(data)
	if length == 0 {
		return 0, nil
	}
	s.length += uint64(length)

	for _, c := range data {

		// TODO: Might want to check this for endianness for consistency
		// across systems.
		if s.offset < 8 {
			s.k1 |= uint64(uint64(c) << uint64(s.offset*8))
		} else if s.offset >= 8 && s.offset < 16 {
			s.k2 |= uint64(uint64(c) << uint64((s.offset%8)*8))
		}
		s.offset++

		if s.offset == 16 {
			s.k1 *= c1_64_128
			s.k1 = (s.k1 << 31) | (s.k1 >> (64 - 31))
			s.k1 *= c2_64_128
			s.h1 ^= s.k1

			s.h1 = (s.h1 << 27) | (s.h1 >> (64 - 27))
			s.h1 += s.h2
			s.h1 = s.h1*5 + 0x52dce729

			s.k2 *= c2_64_128
			s.k2 = (s.k2 << 33) | (s.k2 >> (64 - 33))
			s.k2 *= c1_64_128
			s.h2 ^= s.k2

			s.h2 = (s.h2 << 31) | (s.h2 >> (64 - 31))
			s.h2 += s.h1
			s.h2 = s.h2*5 + 0x38495ab5

			s.k1 = 0
			s.k2 = 0

			s.offset = 0
		}
	}
	return length, nil
}

func (s *sum64_128) Sum128() (uint64, uint64) {
	var h1, h2 uint64 = s.h1, s.h2
	var k1, k2 uint64 = s.k1, s.k2

	//tail
	switch {
	case s.offset > 8:
		k2 *= c2_64_128
		k2 = (k2 << 33) | (k2 >> (64 - 33))
		k2 *= c1_64_128
		h2 ^= k2
		fallthrough

	case s.offset > 0:
		k1 *= c1_64_128
		k1 = (k1 << 31) | (k1 >> (64 - 31))
		k1 *= c2_64_128
		h1 ^= k1
	}

	//finalization
	h1 ^= s.length
	h2 ^= s.length

	h1 += h2
	h2 += h1

	h1 ^= h1 >> 33
	h1 *= 0xff51afd7ed558ccd
	h1 ^= h1 >> 33
	h1 *= 0xc4ceb9fe1a85ec53
	h1 ^= h1 >> 33

	h2 ^= h2 >> 33
	h2 *= 0xff51afd7ed558ccd
	h2 ^= h2 >> 33
	h2 *= 0xc4ceb9fe1a85ec53
	h2 ^= h2 >> 33

	h1 += h2
	h2 += h1

	return h1, h2
}

func (s *sum64_128) Sum(in []byte) []byte {
	h1, h2 := s.Sum128()
	return append(in, byte(h1>>56), byte(h1>>48), byte(h1>>40), byte(h1>>32),
		byte(h1>>24), byte(h1>>16), byte(h1>>8), byte(h1), byte(h2>>56),
		byte(h2>>48), byte(h2>>32), byte(h2>>24), byte(h2>>16),
		byte(h2>>8), byte(h2))
}

func (s *sum64_128) BlockSize() int { return 16 }

func (s *sum64_128) Size() int { return 16 }
