package snefru

import (
	"hash"
)

// reference: https://github.com/mythrill/HashLib/blob/30a1d30e0545db424606e71df367e632b259537d/HashLib/Crypto/Snefru.cs
type snefru struct {
	rounds int
	state  []uint32

	x  []byte // temporary buffer
	i  int    // index into x
	tc uint64 // total count of bytes processed
}

func NewSnefru256(rounds int) hash.Hash {
	h := &snefru{
		rounds: rounds,
		state:  make([]uint32, 8),
	}

	h.x = make([]byte, h.BlockSize())
	h.Reset()

	return h
}

func (h *snefru) _block(b []byte) int {
	nn := 0

	work := make([]uint32, 16)
	copy(work, h.state)
	for len(b) >= h.BlockSize() {
		j := 0
		for i := 0; i < 16; i++ {
			work[i] = uint32(b[j]) | uint32(b[j+1])<<8 | uint32(b[j+2])<<16 | uint32(b[j+3])<<24
			j += 4
		}

		// Rounds
		for i := 0; i < h.rounds; i++ {
			sbox0 := sBoxes[i*2]
			sbox1 := sBoxes[i*2+1]

			for j := 0; j < 4; j++ {
				work[15] ^= sbox0[byte(work[0])]
				work[1] ^= sbox0[byte(work[0])]
				work[0] ^= sbox0[byte(work[1])]
				work[2] ^= sbox0[byte(work[1])]
				work[1] ^= sbox1[byte(work[2])]
				work[3] ^= sbox1[byte(work[2])]
				work[2] ^= sbox1[byte(work[3])]
				work[4] ^= sbox1[byte(work[3])]
				work[3] ^= sbox0[byte(work[4])]
				work[5] ^= sbox0[byte(work[4])]
				work[4] ^= sbox0[byte(work[5])]
				work[6] ^= sbox0[byte(work[5])]
				work[5] ^= sbox1[byte(work[6])]
				work[7] ^= sbox1[byte(work[6])]
				work[6] ^= sbox1[byte(work[7])]
				work[8] ^= sbox1[byte(work[7])]
				work[7] ^= sbox0[byte(work[8])]
				work[9] ^= sbox0[byte(work[8])]
				work[8] ^= sbox0[byte(work[9])]
				work[10] ^= sbox0[byte(work[9])]
				work[9] ^= sbox1[byte(work[10])]
				work[11] ^= sbox1[byte(work[10])]
				work[10] ^= sbox1[byte(work[11])]
				work[12] ^= sbox1[byte(work[11])]
				work[11] ^= sbox0[byte(work[12])]
				work[13] ^= sbox0[byte(work[12])]
				work[12] ^= sbox0[byte(work[13])]
				work[14] ^= sbox0[byte(work[13])]
				work[13] ^= sbox1[byte(work[14])]
				work[15] ^= sbox1[byte(work[14])]
				work[14] ^= sbox1[byte(work[15])]
				work[0] ^= sbox1[byte(work[15])]

				shift := shiftTable[j]
				for n := 0; n < len(work); n++ {
					work[n] = (work[n] >> shift) | (work[n] << (32 - shift))
				}
			}
		}

		h.state[0] ^= work[15]
		h.state[1] ^= work[14]
		h.state[2] ^= work[13]
		h.state[3] ^= work[12]

		// 256
		h.state[4] ^= work[11]
		h.state[5] ^= work[10]
		h.state[6] ^= work[9]
		h.state[7] ^= work[8]

		b = b[h.BlockSize():]
		nn += h.BlockSize()
	}

	return nn
}

func (h *snefru) Write(block []byte) (nn int, err error) {
	nn = len(block)
	h.tc += uint64(nn)

	if h.i > 0 {
		n := len(block)
		if n > h.BlockSize()-h.i {
			n = h.BlockSize() - h.i
		}
		for i := 0; i < n; i++ {
			h.x[h.i+i] = block[i]
		}
		h.i += n
		if h.i == h.BlockSize() {
			h._block(h.x[0:])
			h.i = 0
		}
		block = block[n:]
	}

	n := h._block(block)
	block = block[n:]
	if len(block) > 0 {
		h.i = copy(h.x[:], block)
	}

	return nn, nil
}

func (h *snefru) Sum(in []byte) []byte {
	d := *h

	// Padding.  Add a 1 bit and 0 bits until 56 bytes mod 64.
	tc := d.tc
	var tmp [64]byte
	tmp[0] = 0x80
	if tc%64 < 56 {
		d.Write(tmp[0 : 56-tc%64])
	} else {
		d.Write(tmp[0 : 64+56-tc%64])
	}

	// Length in bits.
	tc <<= 3
	for i := uint(0); i < 8; i++ {
		tmp[i] = byte(tc >> (8 * i))
	}
	d.Write(tmp[0:8])

	if d.i != 0 {
		panic("d.i != 0")
	}

	var digest = make([]byte, h.Size())
	for i, s := range d.state {
		digest[i*4] = byte(s)
		digest[i*4+1] = byte(s >> 8)
		digest[i*4+2] = byte(s >> 16)
		digest[i*4+3] = byte(s >> 24)
	}

	return append(in, digest[:]...)
}

func (h *snefru) Reset() {
	h.state = make([]uint32, len(h.state))
	h.i = 0
	h.tc = 0
}

func (h *snefru) Size() int {
	return 32
}

func (h *snefru) BlockSize() int {
	return 64
}
