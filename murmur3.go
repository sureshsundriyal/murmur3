package murmur3

import "hash"

type Hash128 interface {
	hash.Hash
	Sum128() (uint64, uint64)
}

const (
	c1_64_128 = uint64(0x87c37b91114253d5)
	c2_64_128 = uint64(0x4cf5ad432745937f)
)

type sum64_128 struct {
	h1     uint64
	h2     uint64
	k1     uint64
	k2     uint64
	length uint64
	offset uint8
}

func New64_128() Hash128 {
	return &sum64_128{0, 0, 0, 0, 0, 0}
}

func (s *sum64_128) Reset() {
	s.h1, s.h2, s.k1, s.k2, s.length, s.offset = 0, 0, 0, 0, 0, 0
}

func (s *sum64_128) body() {
	s.k1 *= c1_64_128
	s.k1 = (s.k1 << 31) | (s.k1 >> (64 - 31))
	s.h1 ^= s.k1

	s.h1 = (s.h1 << 27) | (s.h1 >> (64 - 27))
	s.h1 += s.h2
	s.h1 = s.h1*5 + 0x52dce729

	s.k2 *= c2_64_128
	s.k2 = (s.k2 << 33) | (s.k2 >> (64 - 33))
	s.k2 *= c1_64_128
	s.h2 ^= s.k2

	s.h2 = (s.h2 << 31) | (s.h2 >> (64 - 32))
	s.h2 += s.h1
	s.h2 = s.h2*5 + 0x38495ab5
}

func (s *sum64_128) Sum128() (uint64, uint64) {
	var h1, h2 uint64 = s.h1, s.h2

	//tail
	switch {
	case s.offset > 8:
		s.k2 *= c2_64_128
		s.k2 = (s.k2 << 33) | (s.k2 >> (64 - 33))
		s.k2 *= c1_64_128
		h2 ^= s.k2
		fallthrough

	case s.offset > 0:
		s.k1 *= c1_64_128
		s.k1 = (s.k1 << 31) | (s.k1 >> (64 - 31))
		s.k1 *= c2_64_128
		h1 ^= s.k1
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

func (s *sum64_128) Sum(b []byte) []byte {
	return make([]byte, 16)
}

func (s *sum64_128) Write(data []byte) (int, error) {
	length := len(data)
	if length == 0 {
		return 0, nil
	}
	s.length += uint64(length)

	for _, c := range data {

		//TODO: Might want to check this for endianness for consistency
		//across systems.
		if s.offset < 8 {
			s.k1 |= uint64(uint64(c) << uint64(s.offset*8))
			s.offset++
		} else if s.offset >= 8 && s.offset < 16 {
			s.k2 |= uint64(uint64(c) << uint64((s.offset%8)*8))
			s.offset++
		} else {
			//something wrong
		}

		if s.offset == 16 {
			s.body()
			s.k1 = 0
			s.k2 = 0
			s.offset = 0
		}
	}
	return length, nil
}

func (s *sum64_128) BlockSize() int { return 16 }

func (s *sum64_128) Size() int { return 16 }
