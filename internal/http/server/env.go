package server

type ServerEnv struct {
	Port string
}

func NewServerEnv() *ServerEnv {
	return &ServerEnv{}
}
