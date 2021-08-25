package utils

func CalPageOffset(page, limit int64) int64 {
	if page <= 0 || limit <= 0 {
		return 0
	}
	offset := (page - 1) * limit
	return int64(offset)
}
