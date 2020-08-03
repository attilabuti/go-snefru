package snefru

type Buffer struct {
	data []byte
	pos  int
}

func NewBuffer(aLength int) *Buffer {
	return &Buffer{
		data: make([]byte, aLength),
		pos:  0,
	}
}

func (b *Buffer) GetBytes() []byte {
	b.pos = 0
	return b.data
}

//func (b *Buffer) GetBytesZeroPadded() []byte {
//	Array.Clear(m_data, m_pos, m_data.Length - m_pos);
//
//	b.pos = 0
//	return b.data
//}

func (b *Buffer) Feed(aData []byte, aStartIndex, aLength *int, aProcessedBytes *uint64) bool {
	if len(aData) == 0 {
		return false
	}

	if *aLength == 0 {
		return false
	}

	length := len(b.data) - b.pos
	if length > *aLength {
		length = *aLength
	}

	copy(b.data[b.pos:], aData[*aStartIndex:*aStartIndex+length])

	b.pos += length

	*aStartIndex = *aStartIndex + length
	*aLength = *aLength - length
	*aProcessedBytes = *aProcessedBytes + uint64(length)

	return b.pos == len(b.data) //IsFULL
}

func (b *Buffer) IsEmpty() bool {
	return b.pos == 0
}

func (b *Buffer) Length() int {
	return len(b.data)
}

func (b *Buffer) Pos() int {
	return b.pos
}

func (b *Buffer) Initialize() {
	b.pos = 0
}
