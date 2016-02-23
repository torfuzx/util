package util

func SplitLong(input uint64) (high uint32, low uint32) {
	return uint32(input >> 32), uint32(input & 0xFFFFFFFF)
}

func JoinLong(high uint32, low uint32) uint64 {
	return uint64((uint64(high) << 32) | uint64(low))
}
