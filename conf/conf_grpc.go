package conf

import (
	"github.com/Waitfantasy/unicorn/util"
)

type GRpcConf struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enableTls"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	ServerName string `yaml:"serverName"`
}

func (c *GRpcConf) Init() error {
	if c.Addr == "" {
		if v, err := util.GetEnv("UNICORN_GRPC_ADDR", "string"); err != nil {
			c.Addr = "0.0.0.0:9001"
		} else {
			c.Addr = v.(string)
		}
	}

	if c.EnableTLS == false {
		if v, err := util.GetEnv("UNICORN_GRPC_TLS", "bool"); err == nil {
			c.EnableTLS = v.(bool)
		}
	}

	if c.EnableTLS {
		if c.CertFile == "" {
			if v, err := util.GetEnv("UNICORN_GRPC_CERT_FILE_PATH", "string"); err != nil {
				return err
			} else {
				c.CertFile = v.(string)
			}
		}

		if c.KeyFile == "" {
			if v, err := util.GetEnv("UNICORN_GRPC_KEY_FILE_PATH", "string"); err != nil {
				return err
			} else {
				c.KeyFile = v.(string)
			}
		}

		if c.ServerName == "" {
			if v, err := util.GetEnv("UNICORN_GRPC_SERVER_NAME", "string"); err != nil {
				return err
			} else {
				c.ServerName = v.(string)
			}
		}
	}
	return nil
}
