package id

import (
	"github.com/Waitfantasy/unicorn/conf"
	"sync"
	//"sync"
	"testing"
)


func TestId_Gen(t *testing.T) {
	wg := sync.WaitGroup{}
	id := NewAtomicId(&IdConfig{
		MachineId:   1,
		Version:     0,
		ReleaseType: 1,
		IdGenType:   conf.IdPeakGenType,
		Epoch:       1538473327172,
	})
	m := sync.Map{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			uuid, err := id.Gen()
			if err != nil {
				t.Error(err)
			}
			if _, ok := m.Load(uuid); ok {
				t.Error("test fail")
			} else {
				m.Store(uuid, uuid)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
