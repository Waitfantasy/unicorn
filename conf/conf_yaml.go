package conf

import (
	"fmt"
	"github.com/Waitfantasy/unicorn/service/machine"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlConf struct {
	Id            *IdConfig   `yaml:"id"`
	Etcd          *EtcdConfig `yaml:"etcd"`
	Http          *HttpConfig `json:"http"`
	GRpc          *RpcConfig  `yaml:"grpc"`
	Log           *LogConfig  `yaml:"log"`
}

func ParseConfigData(data []byte) (*YamlConf, error) {
	config := &YamlConf{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func ParseConfigFile(filename string) (*YamlConf, error) {
	if filename == "" {
		return &YamlConf{
			Id:   new(IdConfig),
			Http: new(HttpConfig),
			Etcd: new(EtcdConfig),
			GRpc: new(RpcConfig),
			Log:  new(LogConfig),
		}, nil
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfigData(data)
}

func (c *YamlConf) fromEnvInitConfig() error {
	if err := c.Id.fromEnvInitConfig(); err != nil {
		return err
	}

	if err := c.Etcd.fromEnvInitConfig(); err != nil {
		return err
	}

	if err := c.Http.fromEnvInitConfig(); err != nil {
		return err
	}

	if err := c.GRpc.fromEnvInitConfig(); err != nil {
		return err
	}

	c.Log.fromEnvInitConfig()

	return nil
}

func (c *YamlConf) fromEtcdGetMachineItem(ip string) (*machine.Item, error) {
	var err error
	var m *machine.EtcdMachine
	var item *machine.Item

	// create machineService
	if m, err = machine.NewEtcdMachine(*c.Etcd.GetClientV3Config(), c.Etcd.Timeout); err != nil {
		return nil, err
	}

	defer m.Close()

	if item, err = m.Get(ip); err != nil {
		return nil, err
	} else if item != nil {
		return item, nil
	}

	// 向etcd 注册一个新的机器节点
	if item, err = m.Put(ip); err != nil {
		return nil, err
	}

	return item, nil
}

func (c *YamlConf) Init() error {
	// 尝试从env中过去配置
	if err := c.fromEnvInitConfig(); err != nil {
		return err
	}

	// 初始化etcd v3 client 配置
	if err := c.Etcd.initClientV3Config(); err != nil {
		return err
	}

	// 初始化机器id
	if c.Id.MachineIdType == MachineIdEtcdType {
		if item, err := c.fromEtcdGetMachineItem(c.Id.MachineIp); err != nil {
			return fmt.Errorf("using this ip: %s to get machine id from etcd error: %v", c.Id.MachineIp, err)
		} else {
			c.Id.MachineId = item.Id
		}
	}

	return nil
}

func (c *YamlConf) GetIdConfig() *IdConfig {
	return c.Id
}

func (c *YamlConf) GetHttpConfig() *HttpConfig {
	return c.Http
}

func (c *YamlConf) GetEtcdConfig() *EtcdConfig {
	return c.Etcd
}

func (c *YamlConf) GetGRpcConfig() *RpcConfig {
	return c.GRpc
}

func (c *YamlConf) GetLogConfig() *LogConfig {
	return c.Log
}
