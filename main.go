// @title Shortly API DOC
// @version 1.0
// @description Shortly is URL Shortener
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.shortly.io/support
// @contact.email support@shortly.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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
