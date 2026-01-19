package main

import (
	"gymlog/adapters/server"
	"log"
)

func main() {
	gymlogServer := server.NewServer()
	log.Fatal(gymlogServer.Start())
}
