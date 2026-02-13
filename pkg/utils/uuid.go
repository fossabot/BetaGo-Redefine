package utils

import (
	"strconv"
	"time"
)

func GenUUIDStr(str string, length int) string {
	st := strconv.Itoa(int(time.Now().Truncate(time.Minute * 2).Unix()))
	str = st + str
	if len(str) > length {
		return str[:length]
	}
	return str
}

func GenUUIDCode(srcKey, specificKey string, length int) string {
	// 重点是分发的时候，是某个消息（srcKey），只会触发一次（specificKey）
	res := srcKey + specificKey
	if len(res) > length {
		res = res[:length]
	}
	return res
}
