package utils

func Align4K(len uint64) int64 {
	if len == 0 {
		return 1 << 12
	}
	return int64((((len - 1) >> 12) + 1) << 12)
}
