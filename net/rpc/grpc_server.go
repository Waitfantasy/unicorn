package rpc

import (
	"context"
	"github.com/Waitfantasy/unicorn/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type TaskServer struct {
	c conf.Confer
}

func NewTaskServer(c conf.Confer) *TaskServer {
	return &TaskServer{
		c: c,
	}
}

func (s *TaskServer) GetUUID(ctx context.Context, void *Void) (*ResponseUUID, error) {
	g := s.c.GetGenerator()
	res := &ResponseUUID{}
	uuid, err := g.Make()
	if err != nil {
		return nil, err
	}

	res.Uuid = uuid
	return res, nil
}

func (s *TaskServer) Extract(ctx context.Context, extract *RequestExtract) (*ResponseExtract, error) {
	g := s.c.GetGenerator()
	data := g.Extract(extract.Uuid)
	return &ResponseExtract{
		MachineId: int64(data.MachineId),
		Sequence:  data.Sequence,
		Timestamp: data.Timestamp,
		Reserved:  int64(data.Reserved),
		IdType:    int64(data.IdType),
		Version:   int64(data.Version),
	}, nil
}

func (s *TaskServer) ListenAndServe() error {
	var grpcServer *grpc.Server

	// create grpc server
	grpcConf := s.c.GetGRpcConf()
	if grpcConf.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(grpcConf.CertFile, grpcConf.KeyFile)
		if err != nil {
			return err
		} else {
			grpcServer = grpc.NewServer(grpc.Creds(creds))
		}
	} else {
		grpcServer = grpc.NewServer()
	}

	// create listen
	if l, err := net.Listen("tcp", grpcConf.Addr); err != nil {
		return err
	} else {
		RegisterTaskServiceServer(grpcServer, s)
		return grpcServer.Serve(l)
	}
}
