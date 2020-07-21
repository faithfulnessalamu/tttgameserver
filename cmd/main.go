package main

import (
	"log"

	flag "github.com/spf13/pflag"
)

//flags
var (
	port string
)

func main() {
	parseFlags()
	log.Printf("main.Main Port: %s", port)
}

func parseFlags() {
	flag.StringVarP(&port, "port", "p", "8080", "port to run server on")

	flag.Parse()
}
