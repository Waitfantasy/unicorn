package server

import (
	"context"
	"github.com/Waitfantasy/unicorn/id"
	"github.com/Waitfantasy/unicorn/rpc"
	"github.com/Waitfantasy/unicorn/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type Server struct {
	cfg       *rpc.Config
	idService *id.AtomicGenerator
}

func New(cfg *rpc.Config, idService *id.AtomicGenerator) *Server {
	return &Server{
		cfg:       cfg,
		idService: idService,
	}
}

func (s *Server) MakeUUID(ctx context.Context, void *pb.Void) (*pb.MakeResponse, error) {
	res := &pb.MakeResponse{}
	uuid, err := s.idService.Make()
	if err != nil {
		return nil, err
	}

	res.Uuid = uuid
	return res, nil
}

func (s *Server) Transfer(ctx context.Context, request *pb.TransferRequest) (*pb.TransferResponse, error) {
	data := s.idService.Extract(request.Uuid)
	return &pb.TransferResponse{
		MachineId: int64(data.MachineId),
		Sequence:  data.Sequence,
		Timestamp: data.Timestamp,
		Reserved:  int64(data.Reserved),
		IdType:    int64(data.IdType),
		Version:   int64(data.Version),
	}, nil
}

func (s *Server) Run() error {
	var grpcServer *grpc.Server

	// create grpc server
	if s.cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFile, s.cfg.KeyFile)
		if err != nil {
			return err
		} else {
			grpcServer = grpc.NewServer(grpc.Creds(creds))
		}
	} else {
		grpcServer = grpc.NewServer()
	}

	// create listen
	if l, err := net.Listen("tcp", s.cfg.Addr); err != nil {
		return err
	} else {
		pb.RegisterTaskServiceServer(grpcServer, s)
		return grpcServer.Serve(l)
	}
}
