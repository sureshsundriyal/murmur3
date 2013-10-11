package murmur3

import "hash"

type Hash128 interface {
	hash.Hash
	Sum128() (uint64, uint64)
}

const (
	c1_32_128 = uint32(0x239b961b)
	c2_32_128 = uint32(0xab0e9789)
	c3_32_128 = uint32(0x38b34ae5)
	c4_32_128 = uint32(0xa1e38b93)

	c1_64_128 = uint64(0x87c37b91114253d5)
	c2_64_128 = uint64(0x4cf5ad432745937f)
)

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

type sum64_128 struct {
	h1     uint64
	h2     uint64
	k1     uint64
	k2     uint64
	length uint64
	offset uint8
}

func New32_128() Hash128 { return &sum32_128{0, 0, 0, 0, 0, 0, 0, 0, 0, 0} }
func New64_128() Hash128 { return &sum64_128{0, 0, 0, 0, 0, 0} }

func (s *sum32_128) Reset() {
	s.h1, s.h2, s.h3, s.h4 = 0, 0, 0, 0
	s.k1, s.k2, s.k3, s.k4 = 0, 0, 0, 0
	s.length, s.offset = 0, 0
}

func (s *sum64_128) Reset() {
	s.h1, s.h2, s.k1, s.k2, s.length, s.offset = 0, 0, 0, 0, 0, 0
}

func (s *sum32_128) Sum128() (uint64, uint64) {
	var h1, h2, h3, h4 = s.h1, s.h2, s.h3, s.h4

	//tail
	switch {
	case s.offset > 12:
		s.k4 *= c4_32_128
		s.k4 = (s.k4 << 18) | (s.k4 >> (32 - 18))
		s.k4 *= c1_32_128
		h4 ^= s.k4
		fallthrough

	case s.offset > 8:
		s.k3 *= c3_32_128
		s.k3 = (s.k3 << 17) | (s.k3 >> (32 - 17))
		s.k3 *= c4_32_128
		h3 ^= s.k3
		fallthrough

	case s.offset > 4:
		s.k2 *= c2_32_128
		s.k2 = (s.k2 << 16) | (s.k3 >> (32 - 16))
		s.k2 *= c3_32_128
		h2 ^= s.k2
		fallthrough

	case s.offset > 0:
		s.k1 *= c1_32_128
		s.k1 = (s.k1 << 15) | (s.k1 >> (32 - 15))
		s.k1 *= c2_32_128
		h1 ^= s.k1
	}

	//finalization
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

	return uint64((uint64(h1) << 32) | uint64(h2)),
		uint64((uint64(h3) << 32) | uint64(h4))
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

func (s *sum32_128) Sum(in []byte) []byte {
	h1, h2 := s.Sum128()
	return append(in, byte(h1>>56), byte(h1>>48), byte(h1>>40), byte(h1>>32),
		byte(h1>>24), byte(h1>>16), byte(h1>>8), byte(h1), byte(h2>>56),
		byte(h2>>48), byte(h2>>32), byte(h2>>24), byte(h2>>16),
		byte(h2>>8), byte(h2))
}

func (s *sum64_128) Sum(in []byte) []byte {
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

			s.offset = 0
		}
	}
	return length, nil
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

			s.h2 = (s.h2 << 31) | (s.h2 >> (64 - 32))
			s.h2 += s.h1
			s.h2 = s.h2*5 + 0x38495ab5
			s.k1 = 0
			s.k2 = 0

			s.offset = 0
		}
	}
	return length, nil
}

func (s *sum32_128) BlockSize() int { return 16 }
func (s *sum64_128) BlockSize() int { return 16 }

func (s *sum32_128) Size() int { return 16 }
func (s *sum64_128) Size() int { return 16 }
