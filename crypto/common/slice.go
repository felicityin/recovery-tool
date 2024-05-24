package common

// PadToLengthBytesInPlace pad {0, ...} to the front of src if len(src) < length
// output length is equal to the parameter length
func PadToLengthBytesInPlace(src []byte, length int) []byte {
	oriLen := len(src)
	if oriLen < length {
		for i := 0; i < length-oriLen; i++ {
			src = append([]byte{0}, src...)
		}
	}
	return src
}
