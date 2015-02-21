package util

func LeadingZeros(d interface{}) int {
	switch t := d.(type) {
	case uint8:
		return lzs(uint64(t), 8)
	case int8:
		return lzs(uint64(t), 8)
	case int16:
		return lzs(uint64(t), 16)
	case uint16:
		return lzs(uint64(t), 16)
	case int32:
		return lzs(uint64(t), 32)
	case uint32:
		return lzs(uint64(t), 32)
	case int64:
		return lzs(uint64(t), 64)
	case uint64:
		return lzs(uint64(t), 64)
	}
	return -1 // unreachable code
}

func lzs(x uint64, s uint) int {
	u := uint64(1 << (s - 1))
	for i := uint(0); i < s; i++ {
		if (u>>i)&x > 0 {
			return int(i)
		}
	}
	return int(s)
}

func Size(d interface{}) int {
	switch d.(type) {
	case uint8, int8: // also applied to byte
		return 8
	case uint16, int16:
		return 16
	case int32, uint32: // also applied to rune
		return 32
	case int64, uint64:
		return 64
	default:
		return -1
	}
}
