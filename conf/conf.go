package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ReleaseType   int   `json:"release_type"`
	MachineIdType string   `json:"machine_id_type"`
	MachineId     int      `json:"machine_id"`
	MachineIp     string   `json:"machine_ip"`
	Cluster       []string `json:"cluster"`
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
