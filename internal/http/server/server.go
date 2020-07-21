package server

import (
	"context"
	"log"
	"net"
	"net/http"
)

var s http.Server

func Shutdown(ctx context.Context) error {
	log.Println("server.Shutdown shutting down")
	return s.Shutdown(ctx)
}

func Run(ctx context.Context, senv *ServerEnv) {
	s.Addr = net.JoinHostPort("", senv.Port)
	log.Printf("server.Run starting server on %s", s.Addr)

	err := http.ListenAndServe(s.Addr, s.Handler)
	if err != nil {
		log.Fatal("server.Run ", err)
	}
}
