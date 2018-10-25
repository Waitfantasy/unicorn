package restful

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Waitfantasy/unicorn/conf"
	"github.com/Waitfantasy/unicorn/id"
	"github.com/Waitfantasy/unicorn/service/machine"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type Server struct {
	c conf.Confer
	e *gin.Engine
	g *id.AtomicGenerator
}

func NewServer(c conf.Confer) *Server {
	return &Server{
		c: c,
	}
}

func (s *Server) ListenAndServe() error {
	httpConf := s.c.GetHttpConf()
	// TODO
	m, err := machine.NewEtcdMachine(s.c.GetEtcdConf().GetClientConfig())
	if err != nil {
		return err
	}

	defer m.Close()

	api := api{
		m: m,
		g: s.c.GetGenerator(),
	}

	s.e = api.register()
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
				Handler: s.e,
				TLSConfig: &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.RequireAndVerifyClientCert,
				},
			}
			return server.ListenAndServeTLS(httpConf.CertFile, httpConf.KeyFile)
		} else {
			server := &http.Server{
				Addr:    httpConf.Addr,
				Handler: s.e,
			}
			return server.ListenAndServeTLS(httpConf.CertFile, httpConf.KeyFile)
		}
	} else {
		server := &http.Server{
			Addr:    httpConf.Addr,
			Handler: s.e,
		}
		return server.ListenAndServe()
	}
}
