package httpserver

import (
	"fmt"
	"net/http"
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
			ReadTimeout:       15,
			ReadHeaderTimeout: 10,
			WriteTimeout:      15,
			IdleTimeout:       30,
			MaxHeaderBytes:    1024,
			Protocols:         &http2Proto,
		},
	}
}

func (s *Server) Close() error {
	return s.server.Close()
}
