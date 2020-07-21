package server

import "net/http"

type ServerEnv struct {
	Port string
	Handler http.Handler
}

func NewServerEnv() *ServerEnv {
	return &ServerEnv{}
}
