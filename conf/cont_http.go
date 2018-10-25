package conf

import "fmt"

type HttpConf struct {
	Addr       string `yaml:"addr"`
	EnableTLS  bool   `yaml:"enableTls"`
	CaFile     string `yaml:"caFile"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	ClientAuth bool   `yaml:"clientAuth"`
}

func (c *HttpConf) Init() error{
	if err := c.validateEnableTLS(); err != nil {
		return err
	}

	if err := c.validateClientAuth(); err != nil {
		return err
	}

	return nil
}

func (c *HttpConf) validateEnableTLS() error {
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

func (c *HttpConf) validateClientAuth() error {
	if c.ClientAuth {
		if c.CaFile == "" {
			return fmt.Errorf("TLS is enabled and client authentication is enabled, please configure ca file\n")
		}

		if c.CertFile == "" {
			return fmt.Errorf("TLS is enabled and client authentication is enabled, please configure cert file\n")
		}

		if c.KeyFile == "" {
			return fmt.Errorf("TLS is enabled and client authentication is enabled, please configure key file\n")
		}
	}

	return nil
}
