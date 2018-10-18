package conf

type Confer interface {
	Validate() (error)
	InitMachineId() error
	GetIdConf() *IdConf
	GetHttpConf() *HttpConf
	GetEtcdConf() *EtcdConf
	fromLocalGetMachineId() (int, error)
	fromEtcdGetMachineId(ip string) (int, error)
}

type Factory struct {
}

func (f Factory) CreateYamlConf(filename string) (Confer, error){
	return InitYamlConf(filename)
}