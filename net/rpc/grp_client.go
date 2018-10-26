package rpc

import (
	"context"
	"github.com/Waitfantasy/unicorn/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TaskClient struct {
	cli *grpc.ClientConn
}

func NewTaskClient(c conf.Confer) (*TaskClient, error) {
	grpcConf := c.GetGRpcConf()
	var (
		err            error
		cli            *grpc.ClientConn
		transportCreds credentials.TransportCredentials
	)
	if grpcConf.EnableTLS {
		if transportCreds, err = credentials.NewClientTLSFromFile(grpcConf.CertFile, grpcConf.ServerName); err != nil {
			return nil, err
		}

		// TODO context
		if cli, err = grpc.Dial(grpcConf.Addr, grpc.WithTransportCredentials(transportCreds)); err != nil {
			return nil, err
		}

		return &TaskClient{
			cli: cli,
		}, nil

	} else {
		if cli, err = grpc.Dial(grpcConf.Addr, grpc.WithInsecure()); err != nil {
			return nil, err
		}

		return &TaskClient{
			cli: cli,
		}, nil
	}
}

func (t *TaskClient) GetUUID() (*ResponseUUID, error) {
	client := NewTaskServiceClient(t.cli)
	return client.GetUUID(context.Background(), &Void{})
}

func (t *TaskClient) Transfer(uuid uint64) (*ResponseExtract, error) {
	client := NewTaskServiceClient(t.cli)
	return client.Extract(context.Background(), &RequestExtract{
		Uuid: uuid,
	})
}

func (t *TaskClient) Close() error {
	return t.cli.Close()
}
