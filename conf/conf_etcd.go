package conf

import "go.etcd.io/etcd/clientv3"

type EtcdConf struct {
	cfg     *clientv3.Config
	Cluster []string `yaml:"cluster"`
}

func (e *EtcdConf) GetClientConfig() clientv3.Config {
	if e.cfg == nil {
		e.cfg = &clientv3.Config{
			Endpoints: e.Cluster,
		}
	}
	return *e.cfg
}
