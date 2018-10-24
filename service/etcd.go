package service

import (
	"context"
	"errors"
	"github.com/Waitfantasy/unicorn/service/machine"
	"github.com/Waitfantasy/unicorn/util/logger"
	"go.etcd.io/etcd/clientv3"
	"time"
)

const maxEndureMs = 5

type Etcd struct {
	ip  string
	cfg clientv3.Config
}

func NewEtcdService(ip string, cfg clientv3.Config) *Etcd {
	return &Etcd{
		ip:  ip,
		cfg: cfg,
	}
}

func (e *Etcd) VerifyMachineTimestamp() error {
	machineService, err := machine.NewEtcdMachine(e.cfg)
	if err != nil {
		return err
	}

	defer machineService.Close()

	item, err := machineService.Get(e.ip)
	if err != nil {
		return err
	}

	now := machine.Timestamp()
	if now < item.LastTimestamp {
		if offset := item.LastTimestamp - now; offset < maxEndureMs {
			time.Sleep(time.Millisecond * time.Duration(offset<<1))
		} else {
			return errors.New("the current clock has an error")
		}
	}
	return nil
}

func (e *Etcd) ReportMachineTimestamp(ctx context.Context, second time.Duration, l *logger.Log) error {
	var (
		done bool
	)

	machineService, err := machine.NewEtcdMachine(e.cfg)
	if err != nil {
		return err
	}

	defer machineService.Close()

	t := time.NewTimer(second)
	for {
		select {
		case <-ctx.Done():
			done = true
			break
		case <-t.C:
			if item, err := machineService.Get(e.ip); err != nil {
				l.Err("use %s get machine error: %v\n", e.ip, err)
				t.Reset(second)
				break
			} else {
				item.LastTimestamp = machine.Timestamp()
				if err = machineService.PutItem(item); err != nil {
					l.Err("update (ip: %s) machine (last_timestamp: %d) error: %v\n",
						e.ip, item.LastTimestamp, err)
					t.Reset(second)
					break
				}
				l.Debug("update (ip: %s) machine (last_timestamp: %d) success\n", e.ip, item.LastTimestamp)
				t.Reset(second)
			}
		default:
		}

		if done {
			break
		}
	}

	return nil
}
