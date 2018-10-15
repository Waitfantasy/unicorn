package id

import (
	"time"
)

func ExtractTimestamp(idType int, timestamp uint64, epoch uint64) uint64 {
	if idType == SecondIdType {
		return timestamp + (epoch / 1000)
	}

	if idType == MilliSecondIdType {
		return (timestamp + epoch) / 1000
	}

	return timestamp + (epoch / 1000)
}

func Timestamp(idType int, epoch uint64) uint64 {
	// 最大峰值类型, 使用秒级时间戳
	if idType == SecondIdType {
		return (uint64(time.Now().UnixNano()/1000000) - epoch) / 1000
	}

	// 最小粒度类型, 使用毫秒级时间戳
	if idType == MilliSecondIdType {
		return uint64(time.Now().UnixNano()/1000000) - epoch
	}

	return (uint64(time.Now().UnixNano()/1000000) - epoch) / 1000
}

func WaitNextClock(idType int, epoch uint64, lastTimestamp uint64) uint64 {
	ts := uint64(time.Now().UnixNano() / 1000 / 1000)
	for {
		if ts <= lastTimestamp {
			ts = Timestamp(idType, epoch)
		} else {
			break
		}
	}
	return ts
}
