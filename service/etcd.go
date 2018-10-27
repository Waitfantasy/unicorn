package service

import (
	"context"
	"errors"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/service/machine"
	"math"
	"time"
)

const maxEndureMs = 5

type Etcd struct {
	c                   conf.Confer
	m                   *machine.EtcdMachine
	reconnectResultChan chan bool
}

func NewEtcdService(c conf.Confer) (*Etcd, error) {
	if m, err := machine.NewEtcdMachine(c.GetEtcdConf().GetClientConfig(), c.GetEtcdConf().Timeout); err != nil {
		return nil, err
	} else {
		return &Etcd{
			c:                   c,
			m:                   m,
			reconnectResultChan: make(chan bool),
		}, nil
	}
}

func (e *Etcd) VerifyMachineTimestamp() error {
	item, err := e.m.Get(e.c.GetIdConf().MachineIp)
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
		local         = false
		etcd          = true
		retry         = true
		retries       = 0
		maxRetry      = 5
		maxWaitSecond = 10
		done          bool
	)

	second := time.Duration(e.c.GetEtcdConf().Report) * time.Second
	t := time.NewTimer(second)
	l := e.c.GetLogger()
	for {
		select {
		case <-ctx.Done():
			done = true
			break
		case <-t.C:
			// put to local
			if local {
				select {
				case <-e.reconnectResultChan:
					l.Debug("reconnect etcd success, start report timestamp to etcd\n")
					etcd = true
					local = false
					retry = false
					retries = 0
					t.Reset(second)
					break
				default:
					// TODO local storage
					l.Debug("put timestamp tp local\n")
					t.Reset(second)
					break
				}
			}

			// put to etcd
			if etcd {
				if err := e.report(); err != nil {
					retry = true
					retries++
					if retry && retries < maxRetry {
						d := time.Duration(math.Min(math.Pow(2, float64(retries)), float64(maxWaitSecond)))
						l.Debug("report etcd error, retry: %d, wait second: %d\n", retries, d)
						time.Sleep(d)
						t.Reset(second)
						break
					} else {
						l.Debug("max retry, use local put\n")
						local = true
						etcd = false
						go e.reconnect()
						t.Reset(second)
						break
					}
				} else {
					t.Reset(second)
					break
				}
			}
		}

		if done {
			break
		}
	}

	return nil
}

func (e *Etcd) reconnect() {
	l := e.c.GetLogger()
	l.Debug("start reconnect etcd goroutine\n")
	done := false
	// TODO reconnect time use configure
	t := time.NewTimer(time.Second * 3)
	for {
		select {
		case <-t.C:
			l.Debug("[reconnect goroutine]: reconnect etcd\n")
			if err := e.report(); err == nil {
				done = true
				break
			} else {
				t.Reset(time.Second * 3)
				break
			}
		}
		if done {
			e.reconnectResultChan <- true
			t.Stop()
			l.Debug("[reconnect goroutine]: reconnect etcd success\n")
			break
		}
	}
}

func (e *Etcd) report() error {
	ip := e.c.GetIdConf().MachineIp
	item, err := e.m.Get(ip)
	if err != nil {
		return err
	}

	item.LastTimestamp = machine.Timestamp()
	if err = e.m.PutItem(item); err != nil {
		return err
	}

	return nil
}
