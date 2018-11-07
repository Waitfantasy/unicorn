package conf

import (
	"errors"
	"github.com/Waitfantasy/unicorn/util"
)

type RpcConfig struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enableTls"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	ServerName string `yaml:"serverName"`
}

func (c *RpcConfig) fromEnvInitConfig() error {
	if c.Addr == "" {
		if v, err := util.Getenv("UNICORN_GRPC_ADDR", "string"); err != nil {
			c.Addr = "0.0.0.0:6002"
		} else {
			c.Addr = v.(string)
		}
	}

	if c.EnableTLS == false {
		if v, err := util.Getenv("UNICORN_GRPC_TLS", "bool"); err == nil {
			c.EnableTLS = v.(bool)
		}
	}

	if c.EnableTLS {
		if c.CertFile == "" {
			if v, err := util.Getenv("UNICORN_GRPC_CERT_FILE_PATH", "string"); err == nil {
				c.CertFile = v.(string)
			} else {
				return errors.New("grpc service enable tls, but cert file is empty")
			}
		}

		if c.KeyFile == "" {
			if v, err := util.Getenv("UNICORN_GRPC_KEY_FILE_PATH", "string"); err == nil {
				c.KeyFile = v.(string)
			} else {
				return errors.New("grpc service enable tls, but key file is empty")
			}
		}

		if c.ServerName == "" {
			if v, err := util.Getenv("UNICORN_GRPC_SERVER_NAME", "string"); err == nil {
				c.ServerName = v.(string)
			} else {
				return errors.New("grpc service enable tls, but server name is empty")
			}
		}
	}

	return nil
}
