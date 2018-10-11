package time

import (
	"github.com/Waitfantasy/unicorn/conf"
	"time"
)

func Timestamp(idType int, epoch uint64) uint64 {
	// 最大峰值类型, 使用秒级时间戳
	if idType == conf.IdPeakGenType {
		return (uint64(time.Now().UnixNano())/uint64(time.Millisecond) - epoch) / 1000
	}

	// 最小粒度类型, 使用毫秒级时间戳
	if idType == conf.IdGranularityGenType {
		return uint64(time.Now().UnixNano())/uint64(time.Millisecond) - epoch
	}
	return uint64(time.Now().UnixNano())/uint64(time.Millisecond) - epoch
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
