package id

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

type AtomicGenerator struct {
	meta *Meta
	data *Data
	addr *unsafe.Pointer
}

func NewAtomicGenerator(meta *Meta) *AtomicGenerator {
	gen := &AtomicGenerator{
		meta: meta,
		data: &Data{
			0, 0,
		},
	}
	gen.addr = (*unsafe.Pointer)(unsafe.Pointer(gen.data))
	atomic.StorePointer(gen.addr, unsafe.Pointer(gen.data))
	return gen
}

func (gen *AtomicGenerator) Make() (uint64, error) {
	var sequence, timestamp uint64
	for ; ; {
		oldPointer := atomic.LoadPointer(gen.addr)
		oldData := (*Data)(oldPointer)
		timestamp = Timestamp(gen.meta.data.idType, gen.meta.data.epoch)
		sequence = oldData.sequence
		if timestamp < oldData.lastTimestamp {
			return 0, errors.New("clock error.")
		}
		if timestamp == oldData.lastTimestamp {
			sequence++
			sequence &= uint64(gen.meta.GetMaxSequence())
			if sequence == 0 {
				timestamp = WaitNextClock(gen.meta.data.idType, gen.meta.data.epoch, oldData.lastTimestamp)
			}
		} else {
			sequence = 0
		}
		newData := &Data{
			sequence:      sequence,
			lastTimestamp: timestamp,
		}
		if atomic.CompareAndSwapPointer(gen.addr, oldPointer, unsafe.Pointer(newData)) {
			uuid := Calculate(sequence, timestamp, gen.meta)
			return uuid, nil
		}
	}
}

func (gen *AtomicGenerator) Extract(uuid uint64) (*MetaData) {
	data := &MetaData{}
	data.machineId = int(uuid & uint64(gen.meta.GetMaxMachine()))
	data.seq = int((uuid >> gen.meta.GetSequenceLeftShift()) &  uint64(gen.meta.GetMaxSequence()))
	t := uuid >> gen.meta.GetTimestampLeftShift() & uint64(gen.meta.GetMaxTimestamp())
	t = (t * 1000  + gen.meta.data.epoch) / 1000
	data.timestamp = int(t)
	println(t)
	data.service = int((uuid >> gen.meta.GetServiceLeftShift()) & uint64(gen.meta.GetMaxService()))
	data.idType = int((uuid >> gen.meta.GetIdTypeLeftShift()) & uint64(gen.meta.GetMaxIdType()))
	data.version = int((uuid >> gen.meta.GetVersionLeftShift()) & uint64(gen.meta.GetMaxVersion()))
	return data
}
