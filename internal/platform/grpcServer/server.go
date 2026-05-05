package grpcServer

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Config struct {
	Host        string        `env:"HOST" env-prefix:"GRPC_SERVER_" env-required`
	Port        string        `env:"PORT" env-prefix:"GRPC_SERVER_" env-required`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-prefix:"GRPC_SERVER_" env-required`
}

type Server struct {
	lis    net.Listener
	Server *grpc.Server
}

func NewServer(c Config) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", c.Port))
	if err != nil {
		return nil, err
	}
	defer lis.Close()

	server := grpc.NewServer()

	return &Server{
		lis:    lis,
		Server: server,
	}, nil
}

func (s *Server) Serve() error {
	return s.Server.Serve(s.lis)
}

func (s *Server) Close() error {
	return s.lis.Close()
}
