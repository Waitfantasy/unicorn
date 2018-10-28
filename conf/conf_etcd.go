package conf

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
)

const (
	defaultReportSecond      = 30
	defaultTimeoutSecond     = 5
	defaultLocalReportSecond = 90
	defaultReportFile        = "/tmp/unicorn-report.bin"
)

type EtcdConf struct {
	cfg         *clientv3.Config
	EnableTls   bool     `yaml:"enableTls"`
	Insecure    bool     `yaml:"insecure"`
	CaFile      string   `yaml:"caFile"`
	ClientAuth  bool     `yaml:"clientAuth"`
	CertFile    string   `yaml:"certFile"`
	KeyFile     string   `yaml:"keyFile"`
	Report      int      `yaml:"report"`
	Timeout     int      `yaml:"timeout"`
	Cluster     []string `yaml:"cluster"`
	ReportFile  string   `yaml:"reportFile"`
	LocalReport int      `yaml:"localReport"`
}

func (e *EtcdConf) Init() error {
	if err := e.validateEnableTls(); err != nil {
		return err
	}

	if err := e.validateClientAuth(); err != nil {
		return err
	}

	// init report
	if e.Report == 0 {
		e.Report = defaultReportSecond
	}

	// init timeout
	if e.Timeout == 0 {
		e.Timeout = defaultTimeoutSecond
	}

	if e.LocalReport == 0 {
		e.LocalReport = defaultLocalReportSecond
	}

	if e.ReportFile == "" {
		e.ReportFile = defaultReportFile
	}

	// init etcd v3 client
	e.cfg = &clientv3.Config{
		Endpoints: e.Cluster,
	}

	if tlsCfg, err := e.createTlsConfig(); err != nil {
		return err
	} else if tlsCfg != nil {
		e.cfg.TLS = tlsCfg
	}

	return nil
}

func (e *EtcdConf) validateEnableTls() error {
	if e.EnableTls && !e.Insecure {
		if e.CaFile == "" {
			return fmt.Errorf("etcd client TLS is enabled, please configure ca file\n")
		}
	}

	return nil
}

func (e *EtcdConf) validateClientAuth() error {
	if e.ClientAuth {
		if e.CaFile == "" {
			return fmt.Errorf("etcd client enable client auth, please configure ca file\n")
		}

		if e.CertFile == "" {
			return fmt.Errorf("etcd client enable client auth, please configure cert file\n")
		}

		if e.KeyFile == "" {
			return fmt.Errorf("etcd client enable client auth, please configure key file\n")
		}
	}

	return nil
}

func (e *EtcdConf) createTlsConfig() (*tls.Config, error) {
	if e.EnableTls && e.Insecure {
		return &tls.Config{
			InsecureSkipVerify: true,
		}, nil
	}

	var (
		err      error
		certPool *x509.CertPool
	)

	if e.EnableTls {
		if certPool, err = e.createCertPool(); err != nil {
			return nil, err
		}

		return &tls.Config{RootCAs: certPool}, nil
	}

	var certificate tls.Certificate

	if e.ClientAuth {
		if certPool, err = e.createCertPool(); err != nil {
			return nil, err
		}

		if certificate, err = tls.LoadX509KeyPair(e.CertFile, e.KeyFile); err != nil {
			return nil, err
		}

		return &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{certificate},
		}, nil
	}

	return nil, nil
}

func (e *EtcdConf) createCertPool() (*x509.CertPool, error) {
	var (
		err     error
		caBytes []byte
	)

	if caBytes, err = ioutil.ReadFile(e.CaFile); err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caBytes); ok {
		return nil, errors.New("the e clinet use ca file cannot certPool.AppendCertsFromPEM")
	}

	return certPool, nil
}

func (e *EtcdConf) GetClientConfig() clientv3.Config {
	return *e.cfg
}
