package main

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/mux"
	"github.com/sethvargo/go-signalcontext"
	flag "github.com/spf13/pflag"
	"github.com/thealamu/tttgameserver/internal/http/handler"
	"github.com/thealamu/tttgameserver/internal/http/server"

	"github.com/patrickmn/go-cache"
)

//flags
var (
	port string
)

var db = cache.New(cache.NoExpiration, cache.NoExpiration)

func main() {
	parseFlags()
	log.Printf("main.Main port %s", port)

	serverEnv := server.NewServerEnv()
	serverEnv.Port = port

	router := mux.NewRouter()
	router.HandleFunc("/", handler.HomeHandler)
	router.HandleFunc("/newgame", handler.NewGameHandler(db))
	router.HandleFunc("/joingame", handler.JoinGameHandler(db))
	serverEnv.Handler = router

	ctx, cancel := signalcontext.OnInterrupt()
	defer cancel() //call for cleanup

	go server.Run(ctx, serverEnv)

	<-ctx.Done() //wait for CTRL+C

	//stop the server
	ctx = context.Background()
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() //call for cleanup
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal("main.Main ", err)
	}
}

func parseFlags() {
	flag.StringVarP(&port, "port", "p", "8080", "port to run server on")
	flag.Parse()
}
