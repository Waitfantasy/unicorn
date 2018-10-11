package id

//
//import (
//	"errors"
//	"github.com/Soul-Mate/unicorn/conf"
//	"sync"
//	"time"
//)
//
////const (
////	machineBits        = 10
////	peakSeqBits        = 20
////	secondTsBits       = 30
////	granularitySeqBits = 10
////	milliTsBits        = 40
////	releaseTypeBits    = 2
////	idGenTypeBits      = 1
////	versionBits        = 1
////)
//
type IdConfig struct {
	Epoch       uint64
	Version     int
	IdGenType   int
	ReleaseType int
	MachineId   int
}

type IdData struct {
	seq           uint64
	lastTimestamp uint64
}

//
//type data struct {
//	sequence      uint64
//	lastTimestamp uint64
//}
//
//type Id struct {
//	lock   *sync.Mutex
//	config IdConfig
//	data
//}
//
//func NewId(config IdConfig) *Id {
//	return &Id{
//		lock:   new(sync.Mutex),
//		config: config,
//	}
//}
//
//func (id *Id) Gen() (uint64, error) {
//	id.lock.Lock()
//	defer id.lock.Unlock()
//	ts := id.Timestamp()
//	if ts == id.lastTimestamp {
//		id.sequence++
//		id.sequence = id.sequence & id.getMaxSequence()
//		if id.sequence == 0 {
//			ts = id.WaitNextClock(id.lastTimestamp)
//		}
//	} else {
//		id.sequence = 0
//	}
//
//	if ts < id.lastTimestamp {
//		return 0, errors.New("clock error.")
//	}
//	id.lastTimestamp = ts
//	finalSeqBits := id.finalSeqBits()
//	finalTimestampBits := id.finalTimestampBits()
//	sequenceLeftShift := machineBits
//	timestampLeftShift := finalSeqBits + machineBits
//	releaseTypeLeftShift := finalTimestampBits + finalSeqBits + machineBits
//	idGenLeftShift := releaseTypeBits + finalTimestampBits + finalSeqBits + machineBits
//	versionLeftShift := idGenTypeBits + releaseTypeBits + finalTimestampBits + finalSeqBits + machineBits
//	var uuid uint64
//	uuid |= uint64(id.config.MachineId)
//	uuid |= uint64(id.sequence << uint(sequenceLeftShift))
//	uuid |= uint64(ts << timestampLeftShift)
//	uuid |= uint64(id.config.ReleaseType << releaseTypeLeftShift)
//	uuid |= uint64(id.config.IdGenType << idGenLeftShift)
//	uuid |= uint64(id.config.Version << versionLeftShift)
//	return uuid, nil
//}
//
//func (id *Id) Timestamp() uint64 {
//	// 最大峰值类型, 使用秒级时间戳
//	if id.config.IdGenType == conf.IdPeakGenType {
//		return (uint64(time.Now().Nanosecond())/uint64(time.Millisecond))/1000 - id.config.Epoch
//	}
//
//	// 最小粒度类型, 使用毫秒级时间戳
//	if id.config.IdGenType == conf.IdGranularityGenType {
//		return uint64(time.Now().Nanosecond())/uint64(time.Millisecond) - id.config.Epoch
//	}
//	return uint64(time.Now().Nanosecond())/uint64(time.Millisecond) - id.config.Epoch
//}
//
//func (id *Id) WaitNextClock(lastTimestamp uint64) uint64 {
//	ts := uint64(time.Now().UnixNano() / 1000 / 1000)
//	for {
//		if ts <= lastTimestamp {
//			ts = id.Timestamp()
//		} else {
//			break
//		}
//	}
//	return ts
//}
//
//func (id *Id) getMaxSequence() uint64 {
//	// 最大峰值类型, 使用20位序列号
//	if id.config.IdGenType == conf.IdPeakGenType {
//		return -1 ^ (-1 << peakSeqBits)
//	}
//
//	// 最小粒度类型, 使用10位序列号
//	if id.config.IdGenType == conf.IdGranularityGenType {
//		return -1 ^ (-1 << granularitySeqBits)
//	}
//
//	return -1 ^ (-1 << peakSeqBits)
//}
//
//func (id *Id) finalSeqBits() uint {
//	switch id.config.IdGenType {
//	case conf.IdPeakGenType:
//		return peakSeqBits
//	case conf.IdGranularityGenType:
//		return granularitySeqBits
//	default:
//		return peakSeqBits
//	}
//}
//
//func (id *Id) finalTimestampBits() uint {
//	switch id.config.IdGenType {
//	case conf.IdPeakGenType:
//		return secondTsBits
//	case conf.IdGranularityGenType:
//		return milliTsBits
//	default:
//		return secondTsBits
//	}
//}
//
//type BaseId struct {
//	config IdConfig
//}
//
//func (base *BaseId) Timestamp() uint64 {
//	// 最大峰值类型, 使用秒级时间戳
//	if base.config.IdGenType == conf.IdPeakGenType {
//		return (uint64(time.Now().UnixNano()) / uint64(time.Millisecond) - base.config.Epoch) / 1000
//	}
//
//	// 最小粒度类型, 使用毫秒级时间戳
//	if base.config.IdGenType == conf.IdGranularityGenType {
//		return uint64(time.Now().UnixNano()) / uint64(time.Millisecond) - base.config.Epoch
//	}
//	return uint64(time.Now().UnixNano()) / uint64(time.Millisecond) - base.config.Epoch
//}
//
//func (base *BaseId) WaitNextClock(lastTimestamp uint64) uint64 {
//	ts := uint64(time.Now().UnixNano() / 1000 / 1000)
//	for {
//		if ts <= lastTimestamp {
//			ts = base.Timestamp()
//		} else {
//			break
//		}
//	}
//	return ts
//}
//
//func (base *BaseId) GetMaxSequence() uint64 {
//	// 最大峰值类型, 使用20位序列号
//	if base.config.IdGenType == conf.IdPeakGenType {
//		return -1 ^ (-1 << peakSeqBits)
//	}
//
//	// 最小粒度类型, 使用10位序列号
//	if base.config.IdGenType == conf.IdGranularityGenType {
//		return -1 ^ (-1 << granularitySeqBits)
//	}
//
//	return -1 ^ (-1 << peakSeqBits)
//}
//
//func (base *BaseId) FinalSeqBits() uint {
//	switch base.config.IdGenType {
//	case conf.IdPeakGenType:
//		return peakSeqBits
//	case conf.IdGranularityGenType:
//		return granularitySeqBits
//	default:
//		return peakSeqBits
//	}
//}
//
//func (base *BaseId) FinalTimestampBits() uint {
//	switch base.config.IdGenType {
//	case conf.IdPeakGenType:
//		return secondTsBits
//	case conf.IdGranularityGenType:
//		return milliTsBits
//	default:
//		return secondTsBits
//	}
//}
