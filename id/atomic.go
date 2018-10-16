package id

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

type Data struct {
	sequence      uint64
	lastTimestamp uint64
}

type AtomicGenerator struct {
	id   *Id
	data *Data
	addr *unsafe.Pointer
}

func NewAtomicGenerator(id *Id) *AtomicGenerator {
	gen := &AtomicGenerator{
		id:   id,
		data: &Data{0, 0},
	}
	gen.addr = (*unsafe.Pointer)(unsafe.Pointer(gen.data))
	atomic.StorePointer(gen.addr, unsafe.Pointer(gen.data))
	return gen
}

func (gen *AtomicGenerator) Make() (uint64, error) {
	var sequence, timestamp uint64
	for ; ; {
		oldDataPointer := atomic.LoadPointer(gen.addr)
		oldData := (*Data)(oldDataPointer)
		timestamp = gen.id.TimerUtil.Timestamp()
		sequence = oldData.sequence
		if timestamp < oldData.lastTimestamp {
			return 0, errors.New("clock error")
		}

		if timestamp == oldData.lastTimestamp {
			if sequence = (sequence + 1) & uint64(gen.id.Meta.GetMaxSequence()); sequence == 0 {
				timestamp = gen.id.TimerUtil.WaitNextClock(oldData.lastTimestamp)
			}
		} else {
			sequence = 0
		}
		newData := &Data{
			sequence:      sequence,
			lastTimestamp: timestamp,
		}

		if atomic.CompareAndSwapPointer(gen.addr, oldDataPointer, unsafe.Pointer(newData)) {
			uuid := gen.id.calculate(sequence, timestamp)
			return uuid, nil
		}
	}
}

func (gen *AtomicGenerator) Extract(uuid uint64) (*ExtractData) {
	return gen.id.transfer(uuid)
}
