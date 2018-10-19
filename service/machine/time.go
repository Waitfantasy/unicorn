package machine

import "time"

func Timestamp() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}