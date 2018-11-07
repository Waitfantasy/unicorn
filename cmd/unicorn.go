package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/id"
	"github.com/Waitfantasy/unicorn/logger"
	"github.com/Waitfantasy/unicorn/restful"
	restfulServer "github.com/Waitfantasy/unicorn/restful/server"
	"github.com/Waitfantasy/unicorn/rpc"
	rpcServer "github.com/Waitfantasy/unicorn/rpc/server"
	"github.com/Waitfantasy/unicorn/service"
	"github.com/Waitfantasy/unicorn/service/machine"
	"os"
	"runtime"
)

var configFileName = flag.String("config", os.Getenv("UNICORN_CONF"), "unicorn config file")


var (
	BuildDate    string
	BuildVersion string
)

const banner string = `
 ____ __________  .____________  ________ __________  _______   
|    |   \      \ |   \_   ___ \ \_____  \\______   \ \      \  
|    |   /   |   \|   /    \  \/  /   |   \|       _/ /   |   \ 
|    |  /    |    \   \     \____/    |    \    |   \/    |    \
|______/\____|__  /___|\______  /\_______  /____|_  /\____|__  /
                \/            \/         \/       \/         \/ 
`

func main() {
	var (
		err         error
		cfg         *conf.YamlConf
		log         *logger.Logger
		idGenerator *id.AtomicGenerator
		etcdService *service.Etcd
		etcdMachine *machine.EtcdMachine
	)

	fmt.Print(banner)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	fmt.Printf("Git commit:%s\n", BuildVersion)
	fmt.Printf("Build time:%s\n", BuildDate)

	if cfg, err = conf.ParseConfigFile(*configFileName); err != nil {
		fmt.Printf("parse config file error:%v\n", err.Error())
		return
	}

	if err = cfg.Init(); err != nil {
		fmt.Println(err)
		return
	}

	// init atomic id generator
	idCfg, err := id.NewId(&id.Config{
		Epoch:     cfg.Id.Epoch,
		MachineId: cfg.Id.MachineId,
		IdType:    cfg.Id.IdType,
		Version:   cfg.Id.Version,
	})

	if err != nil {
		fmt.Printf("init atomic generator error: %v", err.Error())
		return
	}

	idGenerator = id.NewAtomicGenerator(idCfg)

	// init logger
	log = logger.New(&logger.Config{
		Level:      cfg.Log.Level,
		Output:     cfg.Log.Output,
		Split:      cfg.Log.Split,
		FilePath:   cfg.Log.FilePath,
		FilePrefix: cfg.Log.FilePrefix,
		FileSuffix: cfg.Log.FileSuffix,
	})

	if err = log.Run(); err != nil {
		fmt.Printf("log run error: %v", err.Error())
		return
	}

	// init etcd machine
	if etcdMachine, err = machine.NewEtcdMachine(*cfg.Etcd.GetClientV3Config(), cfg.Etcd.Timeout); err != nil {
		fmt.Printf("init etcd machine error: %v", err.Error())
	}

	defer etcdMachine.Close()

	// create etcd service
	if etcdService, err = service.NewEtcdService(cfg, etcdMachine); err != nil {
		fmt.Printf("create etcd service error: %v", err.Error())
		return
	}

	// verify machine timestamp
	if err = etcdService.VerifyMachineTimestamp(); err != nil {
		fmt.Printf("verify machine timestamp error: %v", err.Error())
		return
	}

	// start report goroutine
	go etcdService.ReportMachineTimestamp(context.Background())

	// run grpc server
	var rpcSrv *rpcServer.Server
	rpcSrv = rpcServer.New(&rpc.Config{
		Addr:       cfg.GRpc.Addr,
		EnableTLS:  cfg.GRpc.EnableTLS,
		CertFile:   cfg.GRpc.CertFile,
		KeyFile:    cfg.GRpc.KeyFile,
		ServerName: cfg.GRpc.ServerName,
	}, idGenerator)

	go rpcSrv.Run()

	// run restful server
	var restfulSrv *restfulServer.Server
	restfulSrv = restfulServer.New(&restful.Config{
		Addr:       cfg.Http.Addr,
		EnableTLS:  cfg.Http.EnableTLS,
		CaFile:     cfg.Http.CaFile,
		CertFile:   cfg.Http.CertFile,
		KeyFile:    cfg.Http.CertFile,
		ClientAuth: cfg.Http.ClientAuth,
	}, idGenerator, etcdMachine)

	if err = restfulSrv.Run(); err != nil {
		fmt.Printf("restful server run error:%v\n", err.Error())
		return
	}
}
