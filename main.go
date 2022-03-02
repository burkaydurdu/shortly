package main

import (
	"log"

	"github.com/burkaydurdu/shortly/config"
	"github.com/burkaydurdu/shortly/internal/server"
)

func checkFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conf, err := config.New()
	checkFatalError(err)

	conf.Print()

	// Create Shortly Server
	shortlyServer := server.NewServer(conf)

	err = shortlyServer.Start()
	checkFatalError(err)
}
