package machine

import "encoding/json"

const maxFreeSlots = MaxMachine + 1

type slots struct {
	free [maxFreeSlots]int
	last int
}



func newSlots() *slots {
	return &slots{
		last:    0,
		free: [maxFreeSlots]int{0: 1024,},
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
	for i := 1; i < MaxMachine; i++ {
		if s.free[i] == 0 {
			return i
		}
	}
	return 0
}
