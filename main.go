package main

import (
	"gymlog/adapters/application"
	"gymlog/adapters/server"
	"gymlog/adapters/storage"
	"log"
)

func main() {
	storage, err := storage.NewSqliteStorage("gymlog.db")
	if err != nil {
		log.Fatal(err)
	}
	routineRepository := application.NewGymRepository(storage)
	userRepository := application.NewUserRepo(storage)
	gymlogServer := server.NewServer(routineRepository, userRepository)
	log.Fatal(gymlogServer.Start())
}
