package utils

func UniqueInt64Slice(s []int64) []int64 {
	if len(s) <= 1 {
		return s
	}
	m := make(map[int64]struct{}, len(s))
	for _, elem := range s {
		m[elem]= struct{}{}
	}
	res := make([]int64, 0, len(m))
	for k, _ := range m {
		res = append(res, k)
	}
	return res
}
