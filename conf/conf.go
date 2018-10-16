package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type etcdConf struct {
	Cluster []string `yaml:"cluster"`
}

type idConf struct {
	Epoch     uint64 `yaml:"epoch"`
	Version   int    `yaml:"version"`
	IdType    int    `yaml:"idType"`
	MachineId int    `yaml:"machineId"`
	MachineIp string `yaml:"machineIp"`
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
