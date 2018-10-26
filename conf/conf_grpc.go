package conf

import "fmt"

type GRpcConf struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enableTls"`
	Insecure   bool   `yaml:"insecure"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	ServerName string `yaml:"serverName"`
}

func (c *GRpcConf) Init() error {
	if err := c.validateEnableTLS(); err != nil {
		return err
	}

	return nil
}

func (c *GRpcConf) validateEnableTLS() error {
	if c.EnableTLS {
		if c.CertFile == "" {
			return fmt.Errorf("TLS is enabled, please configure certFile\n")
		}

		if c.KeyFile == "" {
			return fmt.Errorf("TLS is enabled, please configure keyFile\n")
		}

		if c.ServerName == "" {
			return fmt.Errorf("TLS is enabled, please configure serverName\n")
		}
	}
	return nil
}
