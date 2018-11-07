package conf

import (
	"errors"
	"github.com/Waitfantasy/unicorn/util"
)

const (
	// Machine Id Get Type
	MachineIdLocalType = 0
	MachineIdEtcdType  = 1
)

type IdConfig struct {
	Epoch         uint64 `yaml:"epoch"`
	MachineId     int    `yaml:"machineId"`
	MachineIp     string `yaml:"machineIp"`
	MachineIdType int    `yaml:"machineIdType"`
	IdType        int    `yaml:"idType"`
	Version       int    `yaml:"version"`
}

func (c *IdConfig) fromEnvInitConfig() error{
	if c.Epoch == 0 {
		if v, err := util.Getenv("UNICORN_EPOCH", "uint64"); err == nil {
			c.Epoch = v.(uint64)
		} else {
			return errors.New("epoch can not be empty")
		}
	}

	if v, err := util.Getenv("UNICORN_MACHINE_ID_TYPE", "int"); err == nil {
		c.MachineIdType = v.(int)
	} else {
		c.MachineIdType = MachineIdEtcdType
	}

	if c.MachineIdType == MachineIdLocalType && c.MachineId == 0 {
		if v, err := util.Getenv("UNICORN_MACHINE_ID", "int"); err == nil {
			c.MachineId = v.(int)
		} else {
			return errors.New("you used the local id type in the configuration, please configure it")
		}
	}

	if c.MachineIp == "" {
		if v, err := util.Getenv("UNICORN_MACHINE_IP", "string"); err == nil {
			c.MachineIp = v.(string)
		} else {
			return errors.New("machine ip can not be empty")
		}
	}

	return nil
}