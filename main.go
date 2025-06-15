package main

import (
	"github.com/tylerdevx/geetest-solver/server"
	"log"
)

func main() {
	serverHost := "0.0.0.0"
	serverPort := "8080"

	if err := server.Start(serverHost, serverPort); err != nil {
		log.Fatal(err)
	}
}
