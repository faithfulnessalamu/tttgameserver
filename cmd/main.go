package main

import (
	"context"
	"log"

	flag "github.com/spf13/pflag"
	"github.com/thealamu/tttgameserver/internal/http/server"
)

//flags
var (
	port string
)

func main() {
	parseFlags()
	log.Printf("main.Main port %s", port)

	serverEnv := server.NewServerEnv()
	serverEnv.Port = port

	ctx := context.Background()
	server.Run(ctx, serverEnv)
}

func parseFlags() {
	flag.StringVarP(&port, "port", "p", "8080", "port to run server on")

	flag.Parse()
}
