package conf

import "fmt"

type GRpcConf struct {
	Addr      string `yaml:"addr"`
	EnableTLS bool   `yaml:"enableTls"`
	CertFile  string `yaml:"certFile"`
	KeyFile   string `yaml:"keyFile"`
}

func (c *GRpcConf) Init()  error {
	if err := c.validateEnableTLS(); err != nil {
		return  err
	}

	return nil
}

func (c *GRpcConf) validateEnableTLS() error {
	if c.EnableTLS {
		if c.CertFile == "" {
			return fmt.Errorf("TLS is enabled, please configure cert file\n")
		}

		if c.KeyFile == "" {
			return fmt.Errorf("TLS is enabled, please configure key file\n")
		}
	}
	return nil
}
