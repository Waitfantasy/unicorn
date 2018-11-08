package id

import (
	"sync/atomic"
	"unsafe"
)

type data struct {
	sequence      uint64
	lastTimestamp uint64
}

type AtomicGenerator struct {
	id   *Id
	data *data
	addr *unsafe.Pointer
}

func NewAtomicGenerator(id *Id) *AtomicGenerator {
	gen := &AtomicGenerator{
		id:   id,
		data: &data{0, 0},
	}
	gen.addr = (*unsafe.Pointer)(unsafe.Pointer(gen.data))
	atomic.StorePointer(gen.addr, unsafe.Pointer(gen.data))
	return gen
}



func (gen *AtomicGenerator) Make() (uint64, error) {
	var sequence, timestamp uint64
	for ; ; {
		oldDataPointer := atomic.LoadPointer(gen.addr)
		oldData := (*data)(oldDataPointer)
		timestamp = gen.id.timerUtil.Timestamp()
		sequence = oldData.sequence
		if timestamp == oldData.lastTimestamp {
			if sequence = (sequence + 1) & uint64(gen.id.meta.GetMaxSequence()); sequence == 0 {
				timestamp = gen.id.timerUtil.WaitNextClock(oldData.lastTimestamp)
			}
		} else {
			sequence = 0
		}

		newData := &data{
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
