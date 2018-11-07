package conf

import (
	"errors"
	"github.com/Waitfantasy/unicorn/util"
)


type HttpConfig struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enableTls"`
	CaFile     string `yaml:"caFile"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	ClientAuth bool   `yaml:"clientAuth"`
}

func (c *HttpConfig) fromEnvInitConfig() error {
	if c.Addr == "" {
		if v, err := util.Getenv("UNICORN_HTTP_ADDR", "string"); err != nil {
			c.Addr = "0.0.0.0:6001"
		} else {
			c.Addr = v.(string)
		}
	}

	if c.EnableTLS == false {
		if v, err := util.Getenv("UNICORN_HTTP_TLS", "bool"); err == nil {
			c.EnableTLS = v.(bool)
		}
	}

	if c.ClientAuth == false {
		if v, err := util.Getenv("UNICORN_HTTP_CLIENT_AUTH", "bool"); err == nil {
			c.ClientAuth = v.(bool)
		}
	}

	if c.EnableTLS || c.ClientAuth {
		if c.CertFile == "" {
			if v, err := util.Getenv("UNICORN_HTTP_CERT_FILE_PATH", "string"); err == nil {
				c.CertFile = v.(string)
			} else {
				return errors.New("http service enable tls, but cert file is empty")
			}
		}

		if c.KeyFile == "" {
			if v, err := util.Getenv("UNICORN_HTTP_KEY_FILE_PATH", "string"); err == nil {
				c.KeyFile = v.(string)
			} else {
				return errors.New("http service enable tls, but key file is empty")
			}
		}

		if c.ClientAuth {
			if c.CaFile == "" {
				if v, err := util.Getenv("UNICORN_HTTP_CA_FILE_PATH", "string"); err == nil {
					c.CaFile = v.(string)
				} else {
					return errors.New("http service enable client auth, but ca file is empty")
				}
			}
		}
	}

	return nil
}