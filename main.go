package main

import (
	"context"
	"flag"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/net/restful"
	"github.com/Waitfantasy/unicorn/net/rpc"
	"github.com/Waitfantasy/unicorn/service"
	"log"
	"runtime"
)

var (
	c        conf.Confer
	filename string
)

func main() {
	var (
		err       error
		ctx       context.Context
	)

	flag.StringVar(&filename, "config", "/etc/unicorn/unicorn.yaml", "")
	flag.Parse()

	if c, err = conf.NewYamlConf(filename); err != nil {
		log.Fatalf("conf.NewYamlConf(%s) error: %v\n", filename, err)
	}

	if err = c.Init(); err != nil {
		log.Fatal(err)
	}

	ctx, _ = context.WithCancel(context.Background())

	// verify machine timestamp
	etcdService, err := service.NewEtcdService(c)
	if err != nil {
		log.Fatal(err)
	}

	defer etcdService.Close()

	if err = etcdService.VerifyMachineTimestamp(); err != nil {
		log.Fatalf("etcdService.VerifyMachineTimestamp() error: %v\n", err)
	}

	// grp server
	for i := 0; i < runtime.GOMAXPROCS(runtime.NumCPU()); i++ {
		go func() {
			rpc.NewTaskServer(c).ListenAndServe()
		}()
	}

	// start report
	go func() {
		if etcdService.ReportMachineTimestamp(ctx); err != nil {
			log.Fatalf("etcdService.ReportMachineTimestamp() error: %v\n", err)
		}
	}()

	// restful server
	restfulServer := restful.NewServer(c)

	if err = restfulServer.ListenAndServe(); err != nil {
		log.Fatalf("restfulServer.ListenAndServe() error: %v\n", err)
	}
}
