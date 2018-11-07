package service

import (
	"context"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/logger"
	"github.com/Waitfantasy/unicorn/service/machine"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"
)

const maxEndureMs = 5

type Etcd struct {
	c             conf.Confer
	m             machine.Machiner
	f             *os.File
	reconnectChan chan bool
}

func NewEtcdService(cfg conf.Confer, m machine.Machiner) (*Etcd, error) {
	e := &Etcd{
		c:             cfg,
		m:m,
		reconnectChan: make(chan bool),
	}


	// 创建本地上传文件
	filename := cfg.GetEtcdConfig().LocalReportFile
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if f, err := os.Create(filename); err != nil {
			return nil, err
		} else {
			e.f = f
		}
	} else if err == nil {
		if f, err := os.OpenFile(filename, os.O_RDWR, 0666); err != nil {
			return nil, err
		} else {
			e.f = f
		}
	} else {
		return nil, err
	}

	return e, nil
}

func (e *Etcd) VerifyMachineTimestamp() error {
	var (
		err  error
		item *machine.Item
	)
	for retries := 0; retries < 3; retries++ {
		if item, err = e.m.Get(e.c.GetIdConfig().MachineIp); err != nil {
			retries++
			time.Sleep(time.Duration(math.Min(math.Pow(2, float64(retries)), 5)) * time.Second)
		} else {
			return waitDoubleMachineTimestamp(item.LastTimestamp)
		}
	}

	// 从etcd验证达到最大次数后, 使用本地文件验证
	if ts, err := e.readReportFile(); err != nil {
		return err
	} else {
		return waitDoubleMachineTimestamp(ts)
	}
}

func waitDoubleMachineTimestamp(lastTimestamp uint64) error {
	now := machine.Timestamp()
	if now < lastTimestamp {
		if offset := lastTimestamp - now; offset < maxEndureMs {
			time.Sleep(time.Millisecond * time.Duration(offset<<1))
		} else {
			return fmt.Errorf("the he last synchronized timestamp(%d) is greater than the current timestamp(%d)\n", lastTimestamp, now)
		}
	}
	return nil
}

func (e *Etcd) ReportMachineTimestamp(ctx context.Context) {
	var (
		err   error
		local bool
		done  bool
		etcd  = true
	)
	defer e.f.Close()

	sec1 := time.Duration(e.c.GetEtcdConfig().ReportSec) * time.Second
	sec2 := time.Duration(e.c.GetEtcdConfig().LocalReportSec) * time.Second
	// t1 control to report timestamp to etcd periodically
	t1 := time.NewTicker(sec1)
	// t2 control to report timestamp to local file periodically
	// when an error occurs in the timestamp reported to etcd,
	// a retry is performed. When the retry exceeds the maximum number (default is 5 times),
	// the timestamp is reported to the local file.
	t2 := time.NewTicker(sec2)
	for {
		select {
		case <-t1.C:
			if local {
				select {
				// when receiving the reconnect success message sent by reconnect goroutine,
				// submit the timestamp to etcd
				case <-e.reconnectChan:
					logger.GlobalLogger.Debug("reconnect etcd success, start report timestamp to etcd\n")
					if err = e.retryReport(machine.Timestamp()); err != nil {
						go e.reconnect()
						break
					} else {
						etcd = true
						local = false
						break
					}

				case <-t2.C:
					ts := machine.Timestamp()
					if err = e.writeReportFile(ts); err != nil {
						logger.GlobalLogger.Error("report timestamp-%d to local file error: %v", ts, err)
					}
					logger.GlobalLogger.Debug("report timestamp-%d to local file success\n", ts)
					break
				}
			}

			if etcd {
				if err = e.retryReport(machine.Timestamp()); err != nil {
					local = true
					etcd = false
					go e.reconnect()
				}
				break
			}

		case <-ctx.Done():
			done = true
			break
		}

		if done {
			t1.Stop()
			t2.Stop()
			break
		}
	}
}

func (e *Etcd) readReportFile() (uint64, error) {
	var ts uint64
	if b, err := ioutil.ReadAll(e.f); err != nil {
		return 0, err
	} else {
		if ts, err = strconv.ParseUint(string(b), 10, 64); err != nil {
			return 0, err
		}

		if _, err = e.f.Seek(0, 0); err != nil {
			return 0, err
		}

		return ts, nil
	}
}

func (e *Etcd) writeReportFile(timestamp uint64) error {
	if _, err := e.f.WriteString(strconv.FormatUint(timestamp, 10)); err != nil {
		return err
	}

	if _, err := e.f.Seek(0, 0); err != nil {
		return err
	}

	return nil
}

func (e *Etcd) retryReport(timestamp uint64) error {
	var (
		err           error
		retry         bool
		retries       int
		maxRetry      = 10
		maxWaitSecond = 16
	)
	for {
		if retry {
			timestamp = machine.Timestamp()
		}

		if err = e.report(timestamp); err != nil {
			retry = true
			retries++
			d := time.Duration(math.Min(math.Pow(2, float64(retries)), float64(maxWaitSecond)))
			logger.GlobalLogger.Error("report timestamp-%d to etcd error: %v. retry count: %d, wait second: %d\n",
				timestamp, err, retries, d)
			time.Sleep(d * time.Second)
		} else {
			break
		}

		if !(retry && retries < maxRetry) {
			logger.GlobalLogger.Info("retry to reach the maximum number of times\n")
			break
		}
	}

	return err
}

func (e *Etcd) report(timestamp uint64) error {
	ip := e.c.GetIdConfig().MachineIp
	item, err := e.m.Get(ip)
	if err != nil {
		return err
	}

	item.LastTimestamp = timestamp
	if err = e.m.PutItem(item); err != nil {
		return err
	}

	logger.GlobalLogger.Debug("report timestamp-%d to etcd success\n", item.LastTimestamp)

	return nil
}

func (e *Etcd) reconnect() {
	logger.GlobalLogger.Info("start reconnect etcd goroutine\n")
	var done bool
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-t.C:
			if err := e.report(machine.Timestamp()); err != nil {
				break
			} else {
				done = true
				break
			}
		}

		if done {
			e.reconnectChan <- true
			t.Stop()
			logger.GlobalLogger.Info("[reconnect goroutine]: reconnect etcd success\n")
			break
		}
	}
}