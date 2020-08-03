package snefru

import (
	"hash"
)

// reference: https://github.com/mythrill/HashLib/blob/30a1d30e0545db424606e71df367e632b259537d/HashLib/Crypto/Snefru.cs
type snefru struct {
	rounds int

	state          []uint32
	processedBytes uint64

	blockSize int
	hashSize  int

	buffer *Buffer
}

//
func NewSnefru128(rounds int) hash.Hash {
	return NewSnefru(rounds, 64-16, 16)
}

func NewSnefru256(rounds int) hash.Hash {
	return NewSnefru(rounds, 32, 32)
}

// hashSize => 32
func NewSnefru(rounds int, blockSize int, hashSize int) hash.Hash {
	h := &snefru{
		rounds: rounds,
		state:  make([]uint32, hashSize/4),

		buffer:         NewBuffer(blockSize),
		processedBytes: 0,

		blockSize: blockSize,
		hashSize:  hashSize,
	}

	h.Reset()

	return h
}

func (h *snefru) transformBytes(aData []byte, aIndex, aLength int) {
	if !(h.buffer.IsEmpty()) {
		ok := h.buffer.Feed(aData, &aIndex, &aLength, &h.processedBytes)
		if ok {
			h.transformBuffer()
		}
	}

	for aLength >= h.buffer.Length() {
		h.processedBytes += uint64(h.buffer.Length())
		h.transformBlock(aData, aIndex)
		l := h.buffer.Length()
		aIndex += l
		aLength -= l
	}

	if aLength > 0 {
		h.buffer.Feed(aData, &aIndex, &aLength, &h.processedBytes)
		//aIndex+=nn
		//h.processedBytes += uint64(nn)
	}

	return
}

func (h *snefru) transformBuffer() {
	h.transformBlock(h.buffer.GetBytes(), 0)
}

func (h *snefru) transformBlock(aData []byte, aIndex int) {
	work := make([]uint32, 16)
	copy(work, h.state)

	ConvertBytesToUIntsSwapOrder(aData, aIndex, h.BlockSize(), work, len(h.state))

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
	if h.hashSize == 32 {
		h.state[4] ^= work[11]
		h.state[5] ^= work[10]
		h.state[6] ^= work[9]
		h.state[7] ^= work[8]
	}
}

func (h *snefru) getResult() []byte {
	return ConvertUIntsToBytesSwapOrder(h.state, 0, -1)
}

func (h *snefru) finish() {
	bits := h.processedBytes * 8
	padIndex := 2*h.blockSize - h.buffer.Pos() - 8

	pad := make([]byte, padIndex+8)

	ConvertULongToBytesSwapOrder(bits, pad, padIndex)
	padIndex += 8

	h.transformBytes(pad, 0, padIndex)
}

func (h *snefru) transformFinal() []byte {
	h.finish()

	result := h.getResult()

	h.Reset()
	return result
}

/* go style */

func (h *snefru) Write(block []byte) (nn int, err error) {
	h.transformBytes(block, 0, len(block))

	return len(block), nil
}

func (h *snefru) Sum(in []byte) []byte {
	result := append(in, h.transformFinal()...)
	//h.Reset()
	return result
}

func (h *snefru) Reset() {
	h.state = make([]uint32, len(h.state))
	h.buffer.Initialize()
	h.processedBytes = 0
}

func (h *snefru) Size() int {
	return h.hashSize
}

func (h *snefru) BlockSize() int {
	return h.blockSize
}
