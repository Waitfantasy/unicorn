package service

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/service/machine"
	"math"
	"os"
	"time"
)

const maxEndureMs = 5

type Etcd struct {
	c   conf.Confer
	m   *machine.EtcdMachine
	f   *os.File
	rec chan bool
}

func NewEtcdService(c conf.Confer) (*Etcd, error) {
	e := &Etcd{
		c:   c,
		rec: make(chan bool),
	}

	m, err := machine.NewEtcdMachine(c.GetEtcdConf().GetClientConfig(), c.GetEtcdConf().Timeout)
	if err != nil {
		return nil, err
	}

	e.m = m

	filename := c.GetEtcdConf().ReportFile
	if _, err = os.Stat(filename); os.IsNotExist(err) {
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
		if item, err = e.m.Get(e.c.GetIdConf().MachineIp); err != nil {
			retries++
			time.Sleep(time.Duration(math.Min(math.Pow(2, float64(retries)), 5)) * time.Second)
		} else {
			return waitDoubleMachineTimestamp(item.LastTimestamp)
		}
	}

	// retry limit, read local report file
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
			return errors.New("the current clock has an error")
		}
	}
	return nil
}

func (e *Etcd) ReportMachineTimestamp(ctx context.Context) error {
	var (
		err   error
		local bool = false
		etcd  bool = true
		done  bool = false
	)

	sec1 := time.Duration(e.c.GetEtcdConf().Report) * time.Second
	sec2 := time.Duration(e.c.GetEtcdConf().LocalReport) * time.Second
	// t1 control to report timestamp to etcd periodically
	t1 := time.NewTimer(sec1)
	// t2 control to report timestamp to local file periodically
	// when an error occurs in the timestamp reported to etcd,
	// a retry is performed. When the retry exceeds the maximum number (default is 5 times),
	// the timestamp is reported to the local file.
	t2 := time.NewTimer(sec2)
	l := e.c.GetLogger()
	for {
		select {
		case <-t1.C:
			if local {
				filename := e.c.GetEtcdConf().ReportFile
				select {
				case <-e.rec:
					l.Debug("reconnect etcd success, start report timestamp to etcd\n")
					// when receiving the reconnect success message sent by reconnect goroutine,
					// submit the timestamp of the local file to etcd
					if ts, err := e.readReportFile(); err != nil {
						l.Err("read %s error: %v\n", filename, err)
						etcd = true
						local = false
						t1.Reset(sec1)
						break
					} else {
						// when report to etcd error, start reconnect goroutine
						l.Debug("synchronize local file timestamp-%d to etcd\n", ts)
						if err = e.retryReport(ts); err != nil {
							go e.reconnect()
							t1.Reset(sec1)
							t2.Reset(sec2)
							l.Err("synchronizing local files to etcd error: %v\n", err)
							break
						} else {
							l.Debug("synchronize local files to etcd success\n")
							etcd = true
							local = false
							t1.Reset(sec1)
							break
						}
					}

				case <-t2.C:
					ts := machine.Timestamp()
					if err = e.writeReportFile(ts); err != nil {
						l.Err("report timestamp-%d to local file error: %v", ts, err)
					}
					l.Debug("report timestamp-%d to local file success\n", ts)
					t2.Reset(sec2)
					t1.Reset(sec1)
					break
				}
			}

			// put to etcd
			if etcd {
				if err = e.retryReport(machine.Timestamp()); err != nil {
					local = true
					etcd = false
					go e.reconnect()
					t2.Reset(sec2)
				}
				t1.Reset(sec1)
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

	return nil
}

func (e *Etcd) readReportFile() (uint64, error) {
	var ts uint64
	if err := binary.Read(e.f, binary.BigEndian, &ts); err != nil {
		return ts, err
	}

	if _, err := e.f.Seek(0, 0); err != nil {
		return ts, err
	}

	return ts, nil
}

func (e *Etcd) writeReportFile(timestamp uint64) error {
	if err := binary.Write(e.f, binary.BigEndian, timestamp); err != nil {
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
		maxRetry      = 5
		maxWaitSecond = 10
	)
	l := e.c.GetLogger()
	for {
		if err = e.report(timestamp); err != nil {
			retry = true
			retries++
			sec := time.Duration(math.Min(math.Pow(2, float64(retries)), float64(maxWaitSecond)))
			l.Err("report timestamp-%d to etcd error: %v. retry count: %d, wait second: %d\n",
				timestamp, err, retries, sec)
			time.Sleep(sec)
		} else {
			break
		}

		if !(retry && retries < maxRetry) {
			l.Info("retry to reach the maximum number of times\n")
			break
		}
	}

	return err
}

func (e *Etcd) report(timestamp uint64) error {
	ip := e.c.GetIdConf().MachineIp
	item, err := e.m.Get(ip)
	if err != nil {
		return err
	}

	item.LastTimestamp = timestamp
	if err = e.m.PutItem(item); err != nil {
		return err
	}

	l := e.c.GetLogger()
	l.Debug("report timestamp-%d to etcd success\n", item.LastTimestamp)

	return nil
}

func (e *Etcd) reconnect() {
	l := e.c.GetLogger()
	l.Debug("start reconnect etcd goroutine\n")
	done := false
	sec := time.Minute
	t := time.NewTimer(sec)
	for {
		select {
		case <-t.C:
			l.Debug("[reconnect goroutine]: reconnect etcd\n")
			if err := e.report(machine.Timestamp()); err == nil {
				done = true
				break
			} else {
				t.Reset(sec)
				break
			}
		}
		if done {
			e.rec <- true
			t.Stop()
			l.Debug("[reconnect goroutine]: reconnect etcd success\n")
			break
		}
	}
}

func (e *Etcd) Close() error {
	if err := e.m.Close(); err != nil {
		return err
	}

	if err := e.f.Close(); err != nil {
		return err
	}

	return nil
}
