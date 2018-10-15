package id

import (
	"errors"
	"github.com/Waitfantasy/unicorn-service/machine"
	"go.etcd.io/etcd/clientv3"
)

type EtcdConfig struct {
	Cluster []string `json:"cluster"`
}

type IdConfig struct {
	Epoch     uint64 `json:"epoch"`
	Version   int    `json:"version"`
	IdGenType int    `json:"id_gen_type"`
	// releases type
	// 1: local 2: http 3: rpc
	ReleaseType         int    `json:"release_type"`
	MachineId           int    `json:"machine_id"`
	MachineIp           string `json:"machine_ip"`
	MachineIdAccessType int    `json:"machine_id_type"`
}

type Config struct {
	Etcd *EtcdConfig `json:"etcd"`
	Id   *IdConfig   `json:"id"`
}

func (c *Config) GetMachineId() (int, error) {
	switch c.Id.MachineIdAccessType {
	case MachineIdLocal:
		// TODO machine.ValidMachineId()
		return c.localMachineId()
	case MachineIdEtcd:
		return c.etcdMachineId()
	default:
		return c.localMachineId()
	}
}

func (c *Config) localMachineId() (int, error) {
	// TODO machine.ValidMachineId()
	if !c.ValidId() {
		return 0, errors.New("invalid machine id")
	}
	return c.Id.MachineId, nil
}

func (c *Config) etcdMachineId() (int, error) {
	cfg := c.CreateEtcdClientv3Config()
	service := machine.NewService(cfg)
	if err := service.EtcdConnection(); err != nil {
		return 0, err
	}
	// TODO defer service.EtcdClose()
	if item, err := service.GetMachineItem(service.MachineKey(c.Id.MachineIp)); err != nil {
		return 0, err
	} else {
		return item.Id, nil
	}
}

func (c *Config) ValidId() bool {
	if c.Id.MachineId < machine.MinId || c.Id.MachineId > machine.MaxId {
		return false
	}

	return true
}

func (c *Config) CreateEtcdClientv3Config() clientv3.Config {
	return clientv3.Config{
		Endpoints: c.Etcd.Cluster,
	}
}
