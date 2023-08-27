package main

import (
	"log"

	"github.com/tvvat/project/internal/handler"
	"github.com/tvvat/project/internal/repository"
)

const ()

func main() {
	repo, err := repository.NewPostgresDB(
		repository.PostgresConfig{
			Host:     "localhost",
			Port:     5431,
			Username: "postgres",
			Password: "1",
			DBName:   "admin_parse",
		},
	)
	if err != nil {
		log.Fatalf("can't connect DB: %s", err)
	}

	repo.Migrate()

	hr := handler.NewHandler(repo, ":8080")
	hr.InitRoutes()

	hr.Run()
}
