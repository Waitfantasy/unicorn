package service

import (
	"context"
	"errors"
	"github.com/Waitfantasy/unicorn/service/machine"
			"time"
	"github.com/Waitfantasy/unicorn/conf"
)

const maxEndureMs = 5

type Etcd struct {
	c conf.Confer
}

func NewEtcdService(c conf.Confer) *Etcd {
	return &Etcd{
		c: c,
	}
}

func (e *Etcd) VerifyMachineTimestamp() error {
	machineService, err := machine.NewEtcdMachine(e.c.GetEtcdConf().GetClientConfig())
	if err != nil {
		return err
	}

	defer machineService.Close()

	item, err := machineService.Get(e.c.GetIdConf().MachineIp)
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

func (e *Etcd) ReportMachineTimestamp(ctx context.Context) error {
	var (
		done bool
	)

	machineService, err := machine.NewEtcdMachine(e.c.GetEtcdConf().GetClientConfig())
	if err != nil {
		return err
	}

	defer machineService.Close()

	ip := e.c.GetIdConf().MachineIp
	l := e.c.GetLogger()
	second := time.Duration(e.c.GetEtcdConf().Report)
	t := time.NewTimer(second)
	for {
		select {
		case <-ctx.Done():
			done = true
			break
		case <-t.C:
			if item, err := machineService.Get(ip); err != nil {
				l.Err("use %s get machine error: %v\n", ip, err)
				t.Reset(second)
				break
			} else {
				item.LastTimestamp = machine.Timestamp()
				if err = machineService.PutItem(item); err != nil {
					l.Err("update (ip: %s) machine (last_timestamp: %d) error: %v\n",
						ip, item.LastTimestamp, err)
					t.Reset(second)
					break
				}
				l.Debug("update (ip: %s) machine (last_timestamp: %d) success\n", ip, item.LastTimestamp)
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
