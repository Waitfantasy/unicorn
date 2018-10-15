package id

import (
	"errors"
	"sync"
)

type MutexGenerator struct {
	meta  *Meta
	data  *Data
	mutex sync.Mutex
}

func NewMutexGenerator(meta *Meta) *MutexGenerator {
	return &MutexGenerator{
		meta: meta,
		data: &Data{
			0, 0,
		},
		mutex: sync.Mutex{},
	}
}

func (gen *MutexGenerator) Make() (uint64, error) {
	gen.mutex.Lock()
	defer gen.mutex.Unlock()
	var timestamp uint64
	timestamp = Timestamp(gen.meta.data.idType, gen.meta.data.epoch)
	if timestamp < gen.data.lastTimestamp {
		return 0, errors.New("clock error.")
	}
	if timestamp == gen.data.lastTimestamp {
		if gen.data.sequence = (gen.data.sequence + 1) & uint64(gen.meta.GetMaxSequence()); gen.data.sequence == 0 {
			timestamp = WaitNextClock(gen.meta.data.idType, gen.meta.data.epoch, gen.data.lastTimestamp)
		}
	} else {
		gen.data.sequence = 0
	}

	gen.data.lastTimestamp = timestamp
	uuid := Calculate(gen.data.sequence, gen.data.lastTimestamp, gen.meta)
	return uuid, nil
}

func (gen *MutexGenerator) Extract(uuid uint64) (*MetaData)  {
	panic("impl me")
}
