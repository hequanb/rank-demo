package rediselem

import "strconv"

const NoExpiration = 0

const (
	UserPrefix     = "u:"
	GiftPrefix     = "g:"
	RankingZSetKey = "ranking"
)

func BuildKeys(prefix string, keys ...string) []string {
	res := make([]string, 0, len(keys))
	if len(keys) <= 0 {
		return res
	}

	for _, key := range keys {
		res = append(res, prefix+key)
	}
	return res
}

func BuildKeysInt64(prefix string, keys ...int64) []string {
	res := make([]string, 0, len(keys))
	if len(keys) <= 0 {
		return res
	}

	for _, key := range keys {
		res = append(res, prefix+strconv.FormatInt(key, 10))
	}
	return res
}

func BuildKey(prefix string, key string) string {
	return prefix + key
}
