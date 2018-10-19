package conf

import "go.etcd.io/etcd/clientv3"


type EtcdConf struct {
	Cluster []string `yaml:"cluster"`
}


func (e *EtcdConf) CreateClientV3Config() clientv3.Config {
	return clientv3.Config{
		Endpoints: e.Cluster,
	}
}
