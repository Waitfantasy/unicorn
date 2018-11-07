package client

import (
	"context"
	"github.com/Waitfantasy/unicorn/rpc"
	"github.com/Waitfantasy/unicorn/rpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	cli *grpc.ClientConn
}

func New(c *rpc.Config) (*Client, error) {
	var (
		err            error
		cli            *grpc.ClientConn
		transportCreds credentials.TransportCredentials
	)
	if c.EnableTLS {
		if transportCreds, err = credentials.NewClientTLSFromFile(c.CertFile, c.ServerName); err != nil {
			return nil, err
		}

		// TODO context
		if cli, err = grpc.Dial(c.Addr, grpc.WithTransportCredentials(transportCreds)); err != nil {
			return nil, err
		}

		return &Client{
			cli: cli,
		}, nil

	} else {
		if cli, err = grpc.Dial(c.Addr, grpc.WithInsecure()); err != nil {
			return nil, err
		}

		return &Client{
			cli: cli,
		}, nil
	}
}

func (t *Client) MakeUUID() (*pb.MakeResponse, error) {
	client := pb.NewTaskServiceClient(t.cli)
	return client.MakeUUID(context.Background(), &pb.Void{})
}

func (t *Client) Transfer(uuid uint64) (*pb.TransferResponse, error) {
	client := pb.NewTaskServiceClient(t.cli)
	return client.Transfer(context.Background(), &pb.TransferRequest{
		Uuid: uuid,
	})
}

func (t *Client) Close() error {
	return t.cli.Close()
}
