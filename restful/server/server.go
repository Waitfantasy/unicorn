package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Waitfantasy/unicorn/id"
	"github.com/Waitfantasy/unicorn/restful"
	"github.com/Waitfantasy/unicorn/service/machine"
	"io/ioutil"
	"net/http"
)

type Server struct {
	api *api
	cfg *restful.Config
}

func New(cfg *restful.Config, idService *id.AtomicGenerator, machineService machine.Machiner) *Server {
	return &Server{
		cfg: cfg,
		api: &api{
			idService:      idService,
			machineService: machineService,
		},
	}
}

func (s *Server) Run() error {
	s.api.register()
	if s.cfg.EnableTLS {
		if s.cfg.ClientAuth {
			data, err := ioutil.ReadFile(s.cfg.CaFile)
			if err != nil {
				return fmt.Errorf("failed to read ca certificate: %v\n", err)
			}

			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(data); !ok {
				return fmt.Errorf("add ca certificate failed.\n")
			}

			server := &http.Server{
				Addr:    s.cfg.Addr,
				Handler: s.api.gin,
				TLSConfig: &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.RequireAndVerifyClientCert,
				},
			}
			return server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
		} else {
			server := &http.Server{
				Addr:    s.cfg.Addr,
				Handler: s.api.gin,
			}
			return server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
		}
	} else {
		server := &http.Server{
			Addr:    s.cfg.Addr,
			Handler: s.api.gin,
		}
		return server.ListenAndServe()
	}
}
