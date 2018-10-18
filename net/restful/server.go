package restful

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/id"
	"io/ioutil"
	"net/http"
)

type Server struct {
	c conf.Confer
}

func NewServer(c conf.Confer) *Server {
	return &Server{
		c: c,
	}
}

func (s *Server) ListenAndServe() error {
	httpConf := s.c.GetHttpConf()
	idConf := s.c.GetIdConf()

	handlers := handlers{
		generator: id.NewAtomicGenerator(id.NewId(
			idConf.MachineId,
			idConf.IdType,
			idConf.Version,
			idConf.Epoch)),
	}

	handlers.register()

	if httpConf.EnableTLS {
		if httpConf.ClientAuth {
			data, err := ioutil.ReadFile(httpConf.CaFile)
			if err != nil {
				return fmt.Errorf("failed to read ca certificate: %v\n", err)
			}

			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(data); !ok {
				return fmt.Errorf("add ca certificate failed.\n")
			}

			server := &http.Server{
				Addr:    httpConf.Addr,
				Handler: nil,
				TLSConfig: &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.RequireAndVerifyClientCert,
				},
			}
			return server.ListenAndServeTLS(httpConf.CertFile, httpConf.KeyFile)
		} else {
			return http.ListenAndServeTLS(httpConf.Addr, httpConf.CertFile, httpConf.KeyFile, nil)
		}
	} else {
		return http.ListenAndServe(httpConf.Addr, nil)
	}
}
