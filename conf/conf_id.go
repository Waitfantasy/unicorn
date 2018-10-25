package conf

import (
	"errors"
	"fmt"
	"github.com/Waitfantasy/unicorn/id"
)

type IdConf struct {
	Epoch         uint64 `yaml:"epoch"`
	MachineId     int    `yaml:"machineId"`
	MachineIp     string `yaml:"machineIp"`
	MachineIdType int    `yaml:"machineIdType"`
	IdType        int    `yaml:"idType"`
	Version       int    `yaml:"version"`
}

func (c *IdConf) Init() error {
	if err := c.validateMachineIp(); err != nil {
		return err
	}

	if err := c.validateMachineIdType(); err != nil {
		return err
	}

	if err := c.validateIdType(); err != nil {
		return err
	}

	if err := c.validateVersion(); err != nil {
		return err
	}

	return nil
}

func (c *IdConf) validateMachineIp() error {
	if c.MachineIp == "" {
		return errors.New("Please configure machine ip")
	}
	// TODO regexp validate
	return nil
}

func (c *IdConf) validateMachineIdType() error {
	switch c.MachineIdType {
	case MachineIdLocalType:
		return nil
	case MachineIdEtcdType:
		return nil
	default:
		return fmt.Errorf("the way to get machine id support types: \n\t%d: local type\n\t%d: verify type\n",
			MachineIdLocalType, MachineIdEtcdType)
	}
}

func (c *IdConf) validateIdType() error {
	switch c.IdType {
	case id.SecondIdType:
		return nil
	case id.MilliSecondIdType:
		return nil
	default:
		return fmt.Errorf("id type supports: : \n\t%d: max peak type\n\t%d: min granularity type\n",
			id.SecondIdType, id.MilliSecondIdType)
	}
}

func (c *IdConf) validateVersion() error {
	switch c.Version {
	case UnavailableVersion:
		return nil
	case NormalVersion:
		return nil
	default:
		return fmt.Errorf("version supports: : \n\t%d: max peak type\n\t%d: min granularity type\n",
			UnavailableVersion, NormalVersion)
	}
}

