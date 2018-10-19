package conf

import "github.com/Waitfantasy/unicorn/service/machine"

type Confer interface {
	Validate() error
	InitMachineId() error
	GetIdConf() *IdConf
	GetHttpConf() *HttpConf
	GetEtcdConf() *EtcdConf
	FromLocalGetMachineId() (int, error)
	FromEtcdGetMachineItem(string) (*machine.Item, error)
}

type Factory struct {
}

func (f Factory) CreateYamlConf(filename string) (Confer, error){
	return InitYamlConf(filename)
}