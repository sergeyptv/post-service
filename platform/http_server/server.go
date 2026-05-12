package http_server

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Host string `env:"HOST" env-prefix:"HTTP_" env-required`
	Port string `env:"PORT" env-prefix:"HTTP_" env-required`
}

type Server struct {
	server *http.Server
}

func New(h http.Handler, c Config) *Server {
	http2Proto := http.Protocols{}
	http2Proto.SetUnencryptedHTTP2(true)

	return &Server{
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%s", c.Host, c.Port),
			Handler:           h,
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       30 * time.Second,
			MaxHeaderBytes:    1024,
			Protocols:         &http2Proto,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
