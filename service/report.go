package service

import (
	"fmt"
	"github.com/Waitfantasy/unicorn/service/machine"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func Report(item *machine.Item, second time.Duration, cfg clientv3.Config) error {
	e, err := machine.NewEtcdMachine(cfg)
	if err != nil {
		return err
	}

	defer e.Close()

	t := time.NewTimer(second)
	for {
		select {
		case <-t.C:
			item.LastTimestamp = machine.Timestamp()
			// TODO debug
			fmt.Printf("report timestamp: %d\n", item.LastTimestamp)
			e.PutItem(item)
			t.Reset(second)
		}
	}
}