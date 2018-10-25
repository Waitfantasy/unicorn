package conf

import (
	"github.com/Waitfantasy/unicorn/id"
	"github.com/Waitfantasy/unicorn/util/logger"
)

type Confer interface {
	Init() error
	GetIdConf() *IdConf
	GetHttpConf() *HttpConf
	GetEtcdConf() *EtcdConf
	GetGRpcConf() *GRpcConf
	GetLogConf() *LogConf
	GetGenerator() *id.AtomicGenerator
	GetLogger() *logger.Log
}