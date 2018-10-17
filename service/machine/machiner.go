package machine

const (
	MaxMachine = 1024
	MinMachine = 1
)

type Machiner interface {
	All() ([]*Item, error)
	Get(ip string) (*Item, error)
	Put(ip string) (*Item, error)
	Del(ip string) (*Item, error)
	Replace(id int, ip string) (*Item, *Item, error)
}

func ValidMachineId(id int) bool {
	if id > MaxMachine || id < MinMachine {
		return false
	}

	return true
}
