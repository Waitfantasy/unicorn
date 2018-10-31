package conf

import (
	"github.com/Waitfantasy/unicorn/util"
)


type HttpConf struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enableTls"`
	CaFile     string `yaml:"caFile"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	ClientAuth bool   `yaml:"clientAuth"`
}

func (c *HttpConf) Init() error {
	if c.Addr == "" {
		if v, err := util.GetEnv("UNICORN_HTTP_ADDR", "string"); err != nil {
			c.Addr = "0.0.0.0:6001"
		} else {
			c.Addr = v.(string)
		}
	}

	if c.EnableTLS == false {
		if v, err := util.GetEnv("UNICORN_HTTP_TLS", "bool"); err == nil {
			c.EnableTLS = v.(bool)
		}
	}

	if c.ClientAuth == false {
		if v, err := util.GetEnv("UNICORN_HTTP_CLIENT_AUTH", "bool"); err == nil {
			c.ClientAuth = v.(bool)
		}
	}

	if c.EnableTLS || c.ClientAuth {
		if c.CertFile == "" {
			if v, err := util.GetEnv("UNICORN_HTTP_CERT_FILE_PATH", "string"); err != nil {
				return err
			} else {
				c.CertFile = v.(string)
			}
		}

		if c.KeyFile == "" {
			if v, err := util.GetEnv("UNICORN_HTTP_KEY_FILE_PATH", "string"); err != nil {
				return err
			} else {
				c.KeyFile = v.(string)
			}
		}

		if c.ClientAuth {
			if c.CaFile == "" {
				if v, err := util.GetEnv("UNICORN_HTTP_CA_FILE_PATH", "string"); err != nil {
					return err
				} else {
					c.CaFile = v.(string)
				}
			}
		}
	}
	return nil
}
