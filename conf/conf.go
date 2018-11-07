package conf

type Confer interface {
	Init() error
	GetIdConfig() *IdConfig
	GetEtcdConfig() *EtcdConfig
	GetHttpConfig() *HttpConfig
	GetGRpcConfig() *RpcConfig
	GetLogConfig() *LogConfig
}
