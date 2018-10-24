package main

import (
	"context"
	"flag"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/id"
	"github.com/Waitfantasy/unicorn/net/restful"
	"github.com/Waitfantasy/unicorn/net/rpc"
	"github.com/Waitfantasy/unicorn/service"
	"github.com/Waitfantasy/unicorn/util/logger"
	"go.etcd.io/etcd/clientv3"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	l        *logger.Log
	c        conf.Confer
	filename string
)

func initConf() (err error) {
	factory := conf.Factory{}
	if c, err = factory.CreateYamlConf(filename); err != nil {
		return err
	}

	if err = c.Validate(); err != nil {
		return err
	}

	if err = c.InitMachineId(); err != nil {
		return err
	}

	if l, err = c.GetLogConf().InitLog(); err != nil {
		return err
	}

	return nil
}

func main() {
	var (
		ip        string
		cfg       clientv3.Config
		err       error
		ctx       context.Context
		generator *id.AtomicGenerator
		//cancel context.CancelFunc
	)

	flag.StringVar(&filename, "config", "/etc/unicorn/unicorn.yaml", "")
	flag.Parse()

	if err = initConf(); err != nil {
		log.Fatal(err)
	}

	ip = c.GetIdConf().MachineIp
	cfg = c.GetEtcdConf().GetClientConfig()
	ctx, _ = context.WithCancel(context.Background())
	generator = c.GetIdConf().NewAtomicGenerator()

	// verify machine timestamp
	etcdService := service.NewEtcdService(ip, cfg)
	if err = etcdService.VerifyMachineTimestamp(); err != nil {
		l.Println(logger.LDebug, err)
		os.Exit(1)
	}

	// grp server
	for i := 0; i < runtime.GOMAXPROCS(runtime.NumCPU()); i++ {
		go func() {
			rpc.NewTaskServer(c, generator).ListenAndServe()
		}()
	}

	// start report
	go func() {
		if etcdService.ReportMachineTimestamp(ctx, time.Second*3, l); err != nil {
			l.Err("etcdService.ReportMachineTimestamp error: %v\n", err)
			os.Exit(1)
		}
	}()


	// restful server
	restfulServer := restful.NewServer(c, generator)

	if err = restfulServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
