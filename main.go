package main

import (
	"flag"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/net/restful"
	"log"
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

	if err = c.InitMachineId(); err != nil {
		log.Fatal(err)
	}

	restfulServer := restful.NewServer(c)
	if err = restfulServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
