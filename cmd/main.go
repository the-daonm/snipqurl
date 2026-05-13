package main

import (
	"log"
	"time"

	"snipqurl/internal/database"
	"snipqurl/internal/handler"
	"snipqurl/internal/repository"
	"snipqurl/internal/router"
	"snipqurl/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("no .env file found, using system env")
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal("could not create new db")
	}
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)
	r := router.SetUp(h)

	r.Run()

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			rows, err := repo.DeleteExpired()
			if err != nil {
				log.Printf("cleanup error: %v", err)
				continue
			}
			if rows > 0 {
				log.Printf("cleaned up %d expired urls", rows)
			}
		}
	}()

	r.Run()
}
