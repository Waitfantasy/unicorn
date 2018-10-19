package conf

import (
	"errors"
	"fmt"
	"github.com/Waitfantasy/unicorn/service/machine"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	)

type YamlConf struct {
	Id   *IdConf   `yaml:"id"`
	Http *HttpConf `json:"http"`
	Etcd *EtcdConf `yaml:"verify"`
}

func InitYamlConf(filename string) (*YamlConf, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &YamlConf{}
	if err = yaml.Unmarshal(b, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *YamlConf) Validate() error {
	if err := c.Id.ValidateMachineIp(); err != nil {
		return err
	}

	if err := c.Id.ValidateMachineIdType(); err != nil {
		return err
	}

	if err := c.Id.ValidateIdType(); err != nil {
		return err
	}

	if err := c.Id.ValidateVersion(); err != nil {
		return err
	}

	if err := c.Http.ValidateEnableTLS(); err != nil {
		return err
	}

	if err := c.Http.ValidateClientAuth(); err != nil {
		return err
	}

	return nil
}

func (c *YamlConf) InitMachineId(machiner machine.Machiner) error {
	var (
		id  int
		err error
	)
	switch c.Id.MachineIdType {
	case MachineIdLocalType:
		id, err = c.FromLocalGetMachineId()
	case MachineIdEtcdType:
		item, err := c.FromEtcdGetMachineItem(c.Id.MachineIp, machiner)
		if err != nil {
			return err
		}
		id = item.Id
	default:
		item, err := c.FromEtcdGetMachineItem(c.Id.MachineIp, machiner)
		if err != nil {
			return err
		}
		id = item.Id
	}

	if err != nil {
		return err
	}

	c.Id.MachineId = id
	return nil
}

func (c *YamlConf) FromLocalGetMachineId() (int, error) {
	if !machine.ValidMachineId(c.Id.MachineId) {
		return 0, errors.New("machine id range from 1 ~ 1024")
	}
	return c.Id.MachineId, nil
}

func (c *YamlConf) FromEtcdGetMachineItem(ip string, machiner machine.Machiner) (*machine.Item, error) {
	item, err := machiner.Get(ip)
	if err != nil {
		return nil, err
	}

	if item != nil {
		fmt.Println("from verify get id: ", item.Id)
		return item, nil
	}

	// create new machine id
	item, err = machiner.Put(ip)
	if err != nil {
		return nil, err
	}
	fmt.Println("put verify ip: ", item.Ip)

	return item, nil
}

func (c *YamlConf) GetIdConf() *IdConf {
	return c.Id
}

func (c *YamlConf) GetHttpConf() *HttpConf {
	return c.Http
}

func (c *YamlConf) GetEtcdConf() *EtcdConf {
	return c.Etcd
}
