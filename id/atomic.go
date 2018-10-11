package id

import (
	"errors"
	"github.com/Waitfantasy/unicorn/time"
	"sync/atomic"
	"unsafe"
)

type CASId struct {
	meta     *Meta
	data     *IdData
	dataAddr *unsafe.Pointer
}

func NewAtomicId(c *IdConfig) *CASId {
	casId := new(CASId)
	casId.meta = NewMeta(c)
	casId.data = &IdData{0, 0}
	casId.dataAddr = (*unsafe.Pointer)(unsafe.Pointer(casId.data))
	atomic.StorePointer(casId.dataAddr, unsafe.Pointer(casId.data))
	return casId
}

func (cas *CASId) Gen() (uint64, error) {
	var sequence, ts uint64
	for ; ; {
		oldDataPointer := atomic.LoadPointer(cas.dataAddr)
		oldData := (*IdData)(oldDataPointer)
		ts = time.Timestamp(cas.meta.config.IdGenType, cas.meta.config.Epoch)
		sequence = oldData.seq
		if ts < oldData.lastTimestamp {
			return 0, errors.New("clock error.")
		}

		if ts == oldData.lastTimestamp {
			sequence++
			sequence &= cas.meta.GetMaxSequence()
			if sequence == 0 {
				ts = time.WaitNextClock(cas.meta.config.IdGenType, cas.meta.config.Epoch, oldData.lastTimestamp)
			}
		} else {
			sequence = 0
		}

		newData := IdData{
			seq:           sequence,
			lastTimestamp: ts,
		}

		if atomic.CompareAndSwapPointer(cas.dataAddr, oldDataPointer, unsafe.Pointer(&newData)) {
			uuid := Calculate(sequence, ts, cas.meta)
			return uuid, nil
		}
	}
}
