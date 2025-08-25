package utils

import "time"

func NowTime() int64 {
	return int64(time.Now().Unix())
}