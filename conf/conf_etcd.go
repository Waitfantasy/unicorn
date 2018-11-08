package conf

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/Waitfantasy/unicorn/util"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"strings"
)

const (
	defaultReportSecond      = 3
	defaultTimeoutSecond     = 5
	defaultLocalReportSecond = 6
	defaultReportFile        = "/tmp/unicorn-report.bin"
)

type EtcdConfig struct {
	cfg             *clientv3.Config
	Cluster         []string `yaml:"cluster"`
	EnableTls       bool     `yaml:"enableTls"`
	Insecure        bool     `yaml:"insecure"`
	ClientAuth      bool     `yaml:"clientAuth"`
	CaFile          string   `yaml:"caFile"`
	CertFile        string   `yaml:"certFile"`
	KeyFile         string   `yaml:"keyFile"`
	ReportSec       int      `yaml:"report"`
	Timeout         int      `yaml:"timeout"`
	LocalReportFile string   `yaml:"localReportFile"`
	LocalReportSec  int      `yaml:"localReport"`
}

func (c *EtcdConfig) fromEnvInitConfig() error {
	if c.Cluster == nil {
		if v, err := util.Getenv("UNICORN_ETCD_CLUSTER", "string"); err == nil {
			c.Cluster = strings.Split(v.(string), ";")
		} else {
			return errors.New("etcd cluster can not be empty")
		}
	}

	if c.EnableTls == false {
		if v, err := util.Getenv("UNICORN_ETCD_TLS", "bool"); err == nil {
			c.EnableTls = v.(bool)
		} else {

		}
	}

	if c.Insecure == false {
		if v, err := util.Getenv("UNICORN_ETCD_INSECURE", "bool"); err == nil {
			c.Insecure = v.(bool)
		}
	}

	if c.ClientAuth == false {
		if v, err := util.Getenv("UNICORN_ETCD_CLIENT_AUTH", "bool"); err == nil {
			c.ClientAuth = v.(bool)
		}
	}

	if c.EnableTls {
		if !c.Insecure {
			if c.CaFile == "" {
				if v, err := util.Getenv("UNICORN_ETCD_CA_FILE_PATH", "string"); err == nil {
					c.CaFile = v.(string)
				} else {
					return errors.New("etcd tls is enabled, but ca file is empty")
				}
			}
		}
	}

	if c.ClientAuth {
		if c.CertFile == "" {
			if v, err := util.Getenv("UNICORN_ETCD_CERT_FILE_PATH", "string"); err == nil {
				c.CertFile = v.(string)
			} else {
				return errors.New("etcd client auth is enabled, but cert file is empty")

			}
		}

		if c.KeyFile == "" {
			if v, err := util.Getenv("UNICORN_ETCD_KEY_FILE_PATH", "string"); err == nil {
				c.KeyFile = v.(string)
			} else {
				return errors.New("etcd client auth is enabled, but key file is empty")
			}
		}

		if c.CaFile == "" {
			if v, err := util.Getenv("UNICORN_ETCD_CA_FILE_PATH", "string"); err == nil {
				c.CaFile = v.(string)
			} else {
				return errors.New("etcd client auth is enabled, but ca file is empty")
			}
		}
	}

	if c.Timeout == 0 {
		if v, err := util.Getenv("UNICORN_ETCD_TIMEOUT", "int"); err == nil {
			c.Timeout = v.(int)
		} else {
			c.Timeout = defaultTimeoutSecond
		}
	}

	if c.ReportSec == 0 {
		if v, err := util.Getenv("UNICORN_ETCD_REPORT", "int"); err == nil {
			c.ReportSec = v.(int)
		} else {
			c.ReportSec = defaultReportSecond
		}
	}

	if c.LocalReportSec == 0 {
		if v, err := util.Getenv("UNICORN_ETCD_LOCAL_REPORT", "int"); err == nil {
			c.LocalReportSec = v.(int)
		} else {
			c.LocalReportSec = defaultLocalReportSecond
		}
	}

	if c.LocalReportFile == "" {
		if v, err := util.Getenv("UNICORN_ETCD_LOCAL_REPORT_FILE", "string"); err == nil {
			c.LocalReportFile = v.(string)
		} else {
			c.LocalReportFile = defaultReportFile
		}

	}

	return nil
}

func (c *EtcdConfig) initClientV3Config()  error{
	c.cfg = &clientv3.Config{
		Endpoints: c.Cluster,
	}
	if tlsCfg, err := c.createTlsConfig(); err != nil {
		return err
	} else {
		c.cfg.TLS = tlsCfg
	}

	return nil
}

func (c *EtcdConfig) createCertPool() (*x509.CertPool, error) {
	var (
		err     error
		caBytes []byte
	)

	if caBytes, err = ioutil.ReadFile(c.CaFile); err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caBytes); ok {
		return nil, errors.New("the client use ca file cannot certPool.AppendCertsFromPEM")
	}

	return certPool, nil
}

func (c *EtcdConfig) createTlsConfig() (*tls.Config, error) {
	if c.EnableTls && c.Insecure {
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	var (
		err      error
		certPool *x509.CertPool
	)

	if c.EnableTls {
		if certPool, err = c.createCertPool(); err != nil {
			return nil, err
		}

		return &tls.Config{RootCAs: certPool}, nil
	}

	var certificate tls.Certificate

	if c.ClientAuth {
		if certPool, err = c.createCertPool(); err != nil {
			return nil, err
		}

		if certificate, err = tls.LoadX509KeyPair(c.CertFile, c.KeyFile); err != nil {
			return nil, err
		}

		return &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{certificate},
		}, nil
	}

	return nil, nil
}

func (c *EtcdConfig) GetClientV3Config() *clientv3.Config {
	return c.cfg
}
