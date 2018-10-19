package main

import (
	"flag"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/net/restful"
	"log"
	"github.com/Waitfantasy/unicorn/service/machine"
	"github.com/Waitfantasy/unicorn/service/verify"
)

var filename string

func main() {
	var (
		c conf.Confer
		err error
	)
	flag.StringVar(&filename, "config", "/etc/unicorn/unicorn.yaml", "")
	flag.Parse()

	factory := conf.Factory{}
	if c, err = factory.CreateYamlConf(filename); err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.GetIdConf())
	fmt.Println(c.GetHttpConf())
	fmt.Println(c.GetEtcdConf())

	if err = c.Validate(); err != nil {
		log.Fatal(err)
	}

	// create machiner
	machinerFactory := &machine.MachineFactory{}
	machiner := machinerFactory.CreateEtcdMachine(c.GetEtcdConf().CreateClientV3Config())

	// init machine id
	if err = c.InitMachineId(machiner); err != nil {
		log.Fatal(err)
	}

	// verify machine timestamp
	item, err := c.FromEtcdGetMachineItem(c.GetIdConf().MachineIp, machiner)
	if err != nil {
		log.Fatal(err)
	}

	if err = verify.VerifyMachineTimestamp(item); err != nil {
		log.Fatal(err)
	}

	go machiner.Report(item, 3)

	restfulServer := restful.NewServer(c)
	if err = restfulServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
