package machine

import "go.etcd.io/etcd/clientv3"

const (
	MinMachine = 1
	MaxMachine = 1024
)

type Machiner interface {
	All() ([]*Item, error)
	Get(ip string) (*Item, error)
	Put(ip string) (*Item, error)
	PutItem(item *Item) error
	Del(ip string) (*Item, error)
	Reset(oldIp, newIp string) error
}

type MachineFactory struct {}

func (f MachineFactory) CreateEtcdMachine(config clientv3.Config) (Machiner, error){
	return NewEtcdMachine(config)
}

func ValidMachineId(id int) bool {
	if id > MaxMachine || id < MinMachine {
		return false
	}

	return true
}
