package verify

import (
	"github.com/Waitfantasy/unicorn/service/machine"
	"time"
	"errors"
)

const maxEndureMs = 5

func MachineTimestamp(item *machine.Item) error {
	now := machine.Timestamp()
	if now < item.LastTimestamp {
		if offset := item.LastTimestamp - now; offset < maxEndureMs {
			time.Sleep(time.Millisecond * time.Duration(offset<<1))
		} else {
			return errors.New("the current clock has an error")
		}
	}
	return nil
}

