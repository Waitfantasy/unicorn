package conf

import (
	"errors"
	"github.com/Waitfantasy/unicorn/service/machine"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)


type idConf struct {
	Epoch         uint64 `yaml:"epoch"`
	Version       int    `yaml:"version"`
	IdType        int    `yaml:"idType"`
	MachineId     int    `yaml:"machineId"`
	MachineIp     string `yaml:"machineIp"`
	MachineIdType int    `yaml:"machineIdType"`
}

type Conf struct {
	Etcd *etcdConf `yaml:"etcd"`
	Id   *idConf   `yaml:"id"`
}

func InitConf(filename string) (*Conf, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Conf{}
	if err = yaml.Unmarshal(b, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Conf) GetMachineId() (int, error){
	switch c.Id.MachineIdType {
	case MachineIdLocalType:
		return c.fromLocalGetMachineId()
	case MachineIdEtcdType:
		panic("impl me")
	default:
		panic("impl me")
	}
}

func (c *Conf) fromLocalGetMachineId() (int, error){
	if machine.ValidMachineId(c.Id.MachineId) {
		return 0, errors.New("machine id range from 1 ~ 1024")
	}
	return c.Id.MachineId, nil
}

func (c *Conf) fromEtcdGetMachineId(ip string) (int, error) {
	cfg := c.Etcd.createClientV3Config()

}
