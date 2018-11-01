package main

import (
	"context"
	"flag"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/net/restful"
	"github.com/Waitfantasy/unicorn/net/rpc"
	"github.com/Waitfantasy/unicorn/service"
	"log"
)

var (
	filename string
	c        conf.Confer
)

func main() {
	var (
		err         error
		ctx         context.Context
		etcdService *service.Etcd
	)

	flag.StringVar(&filename, "config", "", "The config file path of the project")
	flag.Parse()

	if c, err = conf.NewYamlConf(filename); err != nil {
		log.Fatalf("conf.NewYamlConf(%s) error: %v\n", filename, err)
	}

	if err = c.Init(); err != nil {
		log.Fatal(err)
	}

	ctx, _ = context.WithCancel(context.Background())

	if etcdService, err = service.NewEtcdService(c); err != nil {
		log.Fatal(err)
	}

	defer etcdService.Close()
	// verify machine timestamp
	if err = etcdService.VerifyMachineTimestamp(); err != nil {
		log.Fatal(err)
	}

	// start log split
	if c.GetLogConf().Split {
		go func() {
			c.GetLogger().Split()
		}()
	}

	// start report
	go etcdService.ReportMachineTimestamp(ctx)

	// grp server
	go func() {
		if err = rpc.NewTaskServer(c).ListenAndServe(); err != nil {
			log.Fatalf("rpc.ListenAndServe() error: %v\n", err)
		}
	}()

	// restful server
	if err = restful.NewServer(c).ListenAndServe(); err != nil {
		log.Fatalf("restful.ListenAndServe() error: %v\n", err)
	}
}
