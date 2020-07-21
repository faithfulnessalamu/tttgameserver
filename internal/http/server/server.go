package server

import (
	"context"
	"log"
	"net"
	"net/http"
)

var s http.Server

func Run(ctx context.Context, senv *ServerEnv) {
	s.Addr = net.JoinHostPort("", senv.Port)
	log.Printf("server.Run starting server on %s", s.Addr)

	err := http.ListenAndServe(s.Addr, s.Handler)
	if err != nil {
		panic(err)
	}
}
