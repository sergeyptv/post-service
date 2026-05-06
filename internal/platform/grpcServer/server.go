package grpcServer

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Config struct {
	Port string `env:"PORT" env-prefix:"GRPC_SERVER_" env-required`
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

	server := grpc.NewServer()

	return &Server{
		lis:    lis,
		Server: server,
	}, nil
}

func (s *Server) Serve() error {
	return s.Server.Serve(s.lis)
}

func (s *Server) CloseListener() error {
	return s.lis.Close()
}

func (s *Server) Shutdown() {
	s.Server.GracefulStop()
}
