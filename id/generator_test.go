package id

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestId_Gen(t *testing.T) {
	wg := sync.WaitGroup{}
	id, err := NewId(&Config{
		Epoch:     1539660973223,
		MachineId: 1,
		IdType:    1,
		Reserved:  0,
		Version:   0,
	})
	if err != nil {
		t.Error(err)
		return
	}

	gen := NewAtomicGenerator(id)
	m := sync.Map{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			uuid, err := gen.Make()
			if err != nil {
				t.Error(err)
			}

			if _, ok := m.Load(uuid); ok {
				t.Error("test fail")
			} else {
				m.Store(uuid, uuid)
				data := gen.Extract(uuid)
				fmt.Printf("machine: %d, seq: %d, timestamp: %s, service: %d, id type: %d, version: %d\n",
					data.MachineId, data.Sequence,
					time.Unix(int64(data.Timestamp), 0).Format("2006-01-02 15:04:05"),
					data.Reserved, data.IdType, data.Version)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
