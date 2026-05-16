package grpc_client

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Config struct {
	Host        string        `env:"HOST" env-required`
	Port        string        `env:"PORT" env-required`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-required`
}

type Client struct {
	Conn *grpc.ClientConn
}

func NewClient(c Config) (*Client, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", c.Host, c.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(c.IdleTimeout),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
