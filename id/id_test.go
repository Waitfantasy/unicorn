package id

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestId_Gen(t *testing.T) {
	wg := sync.WaitGroup{}
	factory := GeneratorFactory{}
	gen := factory.CreateGenerator(AtomicGeneratorType, NewMeta(&MetaData{
		epoch:     1538473327172,
		idType:    SecondIdType,
		service:   1,
		version:   1,
		machineId: 1,
	}))

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
					data.machineId, data.seq, time.Unix(int64(data.timestamp), 0).Format("2006-01-02 15:04:05"), data.service, data.idType, data.version)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
