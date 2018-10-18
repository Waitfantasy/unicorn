package machine

import (
	"encoding/json"
)

const maxFreeSlots = MaxMachine + 1

type slots struct {
	Free [maxFreeSlots]int
	Use  int
}



func newSlots() *slots {
	return &slots{
		Use:  0,
		Free: [maxFreeSlots]int{0: 1024,},
	}
}

func jsonUnmarshalSlots(b []byte) (*slots, error) {
	data := &slots{}
	if err := json.Unmarshal(b, data); err != nil {
		return nil, err
	}
	return data, nil
}

func jsonMarshalSlots(s *slots) ([]byte, error) {
	return json.Marshal(s)
}

func (s *slots) findFreeIndex() int {
	for i := 1; i < maxFreeSlots; i++ {
		if s.Free[i] == 0 {
			return i
		}
	}
	return 0
}
