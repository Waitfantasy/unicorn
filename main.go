package main

import (
	"flag"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/net/restful"
	"github.com/Waitfantasy/unicorn/service"
	"github.com/Waitfantasy/unicorn/service/verify"
	"log"
	"time"
	"os"
)

var filename string

func main() {
	var (
		c   conf.Confer
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

	// init machine id
	if err = c.InitMachineId(); err != nil {
		log.Fatal(err)
	}

	item, err := c.FromEtcdGetMachineItem(c.GetIdConf().MachineIp)
	if err != nil {
		log.Fatal(err)
	}

	// verify machine timestamp
	if err = verify.MachineTimestamp(item); err != nil {
		log.Fatal(err)
	}
	go service.Report(item, time.Second*3, c.GetEtcdConf().CreateClientV3Config())

	restfulServer := restful.NewServer(c)
	if err = restfulServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
