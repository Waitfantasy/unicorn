package conf

import (
	"errors"
	"fmt"
	"github.com/Waitfantasy/unicorn/service/machine"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlConf struct {
	Id     *IdConf   `yaml:"id"`
	Http   *HttpConf `json:"http"`
	Etcd   *EtcdConf `yaml:"etcd"`
	GRpc   *GRpcConf `yaml:"grpc"`
	Log    *LogConf  `yaml:"log"`
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

	if err := c.GRpc.ValidateEnableTLS(); err != nil {
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
		item, err := c.fromEtcdGetMachineItem(c.Id.MachineIp)
		if err != nil {
			return err
		}
		id = item.Id
	default:
		item, err := c.fromEtcdGetMachineItem(c.Id.MachineIp)
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

func (c *YamlConf) fromLocalGetMachineId() (int, error) {
	if !machine.ValidMachineId(c.Id.MachineId) {
		return 0, errors.New("machine id range from 1 ~ 1024")
	}
	return c.Id.MachineId, nil
}

func (c *YamlConf) fromEtcdGetMachineItem(ip string) (*machine.Item, error) {
	// create machineService
	machineService, err := machine.NewEtcdMachine(c.Etcd.GetClientConfig())
	if err != nil {
		return nil, err
	}

	defer machineService.Close()
	item, err := machineService.Get(ip)
	if err != nil {
		return nil, err
	}

	if item != nil {
		fmt.Println("get item: ", item)
		return item, nil
	}

	// create new machine id
	item, err = machineService.Put(ip)
	if err != nil {
		return nil, err
	}
	fmt.Println("put item: ", item)

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

func (c *YamlConf) GetGRpcConf() *GRpcConf {
	return c.GRpc
}

func (c *YamlConf) GetLogConf() *LogConf {
	return c.Log
}

func (c *YamlConf) NewMachine(name string) (machine.Machiner, error) {
	switch name {
	case "etcd":
		cfg := c.GetEtcdConf().GetClientConfig()
		return machine.NewEtcdMachine(cfg)
	default:
		return nil, errors.New("no support machine service")
	}
}
