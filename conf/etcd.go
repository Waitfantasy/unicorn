package conf

import "go.etcd.io/etcd/clientv3"


type etcdConf struct {
	Cluster []string `yaml:"cluster"`
}


func (e *etcdConf) createClientV3Config() *clientv3.Config {
	return &clientv3.Config{
		Endpoints: e.Cluster,
	}
}
