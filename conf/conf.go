package conf

import (
	"encoding/json"
	"io/ioutil"
)

type EtcdConfig struct {
	Cluster []string `json:"cluster"`
}

type IdConfig struct {
	Epoch               uint64 `json:"epoch"`
	Version             int    `json:"version"`
	IdGenType           int    `json:"id_gen_type"`
	ReleaseType         int    `json:"release_type"`
	MachineId           int    `json:"machine_id"`
	MachineIp           string `json:"machine_ip"`
	MachineIdAccessType int    `json:"machine_id_type"`
}

type Config struct {
	Etcd *EtcdConfig `json:"etcd"`
	Id   *IdConfig   `json:"id"`
}

func InitConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err = json.Unmarshal(b, config); err != nil {
		return nil, err
	}
	return config, nil
}
