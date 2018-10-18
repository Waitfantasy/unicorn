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

func (c *HttpConf) ValidateEnableTLS() error {
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

func (c *HttpConf) ValidateClientAuth() error {
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
