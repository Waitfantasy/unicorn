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
	Etcd *EtcdConf `yaml:"etcd"`
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

func (c *YamlConf) InitMachineId() error {
	var (
		id  int
		err error
	)
	switch c.Id.MachineIdType {
	case MachineIdLocalType:
		id, err = c.fromLocalGetMachineId()
	case MachineIdEtcdType:
		id, err = c.fromEtcdGetMachineId(c.Id.MachineIp)
	default:
		id, err = c.fromEtcdGetMachineId(c.Id.MachineIp)
	}

	if err != nil {
		return err
	}

	c.Id.MachineId = id
	return nil
}

func (c *YamlConf) fromLocalGetMachineId() (int, error) {
	if !machine.ValidMachineId(c.Id.MachineId) {
		return 0, errors.New("machine id range from 1 ~ 1024")
	}
	return c.Id.MachineId, nil
}

func (c *YamlConf) fromEtcdGetMachineId(ip string) (int, error) {
	cfg := c.Etcd.createClientV3Config()
	factory := machine.MachineFactory{}
	e := factory.CreateEtcdMachine(cfg)
	if err := e.Conn(); err != nil {
		return 0, err
	}

	defer e.Close()

	item, err := e.Get(ip)
	if err != nil {
		return 0, err
	}

	if item != nil {
		fmt.Println("from etcd get id: ", item.Id)
		return item.Id, nil
	}

	// create new machine id
	item, err = e.Put(ip)
	if err != nil {
		return 0, err
	}
	fmt.Println("put etcd ip: ", item.Ip)

	return item.Id, nil
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
