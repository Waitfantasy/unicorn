package conf

import (
	"testing"
)

func Test_InitConfig(t *testing.T) {
	assertConf := Conf{
		Id: &idConf{
			Epoch:     1539660973223,
			MachineId: 1,
			MachineIp: "127.0.0.1",
			IdType:    0,
			Version:   0,
		},
		Etcd: &EtcdConf{
			Cluster: []string{
				"192.168.10.10:2379",
				"192.168.10.11:2379",
				"192.168.10.12:2379",
			},
		},
	}
	if conf, err := InitConf("./conf.yaml"); err != nil {
		t.Error(err)
	} else {
		for i, v := range conf.Etcd.Cluster {
			if assertConf.Etcd.Cluster[i] != v {
				t.Error("parse conf.Id.Cluster conf.yaml error")
				return
			}
		}

		if conf.Id.Epoch != assertConf.Id.Epoch {
			t.Error("parse conf.Id.Epoch conf.yaml error")
			return
		}

		if conf.Id.MachineId != assertConf.Id.MachineId {
			t.Error("parse conf.Id.MachineId in conf.yaml error")
			return
		}

		if conf.Id.MachineIp != assertConf.Id.MachineIp {
			t.Error("parse conf.Id.MachineIp in conf.yaml error")
			return
		}

		if conf.Id.IdType != assertConf.Id.IdType {
			t.Error("parse conf.Id.IdType in conf.yaml error")
			return
		}

		if conf.Id.Version != assertConf.Id.Version {
			t.Error("parse conf.Id.version in conf.yaml error")
			return
		}
	}
}
