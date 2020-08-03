package snefru

func ConvertBytesToUIntsSwapOrder(aIn []byte, aIndex, aLength int, aResult []uint32, aIndexOut int) {
	for i := aIndexOut; aLength > 0; aLength -= 4 {
		//aIndex++
		aResult[i] = uint32(aIn[aIndex]) << 24
		aIndex++
		aResult[i] |= uint32(aIn[aIndex]) << 16
		aIndex++
		aResult[i] |= uint32(aIn[aIndex]) << 8
		aIndex++
		aResult[i] |= uint32(aIn[aIndex])
		aIndex++

		i++
	}
}

func ConvertULongToBytesSwapOrder(aIn uint64, aOut []byte, aIndex int) {
	aOut[aIndex] = byte(aIn >> 56)
	aIndex++
	aOut[aIndex] = byte(aIn >> 48)
	aIndex++
	aOut[aIndex] = byte(aIn >> 40)
	aIndex++
	aOut[aIndex] = byte(aIn >> 32)
	aIndex++
	aOut[aIndex] = byte(aIn >> 24)
	aIndex++
	aOut[aIndex] = byte(aIn >> 16)
	aIndex++
	aOut[aIndex] = byte(aIn >> 8)
	aIndex++
	aOut[aIndex] = byte(aIn)
	aIndex++
}

func ConvertUIntsToBytesSwapOrder(aIn []uint32, aIndex, aLength int) []byte {
	if aLength == -1 {
		aLength = len(aIn)
	}

	result := make([]byte, aLength*4)

	for j := 0; aLength > 0; {
		result[j] = byte(aIn[aIndex] >> 24)
		j++
		result[j] = byte(aIn[aIndex] >> 16)
		j++
		result[j] = byte(aIn[aIndex] >> 8)
		j++
		result[j] = byte(aIn[aIndex])
		j++

		aLength--
		aIndex++
	}

	return result
}
