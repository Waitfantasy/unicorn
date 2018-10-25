package conf

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"time"
)

const (
	defaultReportSecond = 30
)

type EtcdConf struct {
	cfg        *clientv3.Config
	EnableTls  bool     `yaml:"enableTls"`
	Insecure   bool     `yaml:"insecure"`
	CaFile     string   `yaml:"caFile"`
	ClientAuth bool     `yaml:"clientAuth"`
	CertFile   string   `yaml:"certFile"`
	KeyFile    string   `yaml:"keyFile"`
	Report     int      `yaml:"report"`
	Cluster    []string `yaml:"cluster"`
}

func (e *EtcdConf) Init() error {
	// init report
	if e.Report == 0 {
		e.Report = defaultReportSecond * int(time.Second)
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
