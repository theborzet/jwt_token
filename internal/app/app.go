package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/theborzet/time-tracker/internal/config"
	database "github.com/theborzet/time-tracker/internal/db"
	"github.com/theborzet/time-tracker/internal/delivery/http/handler"
	"github.com/theborzet/time-tracker/internal/delivery/http/routes"
	"github.com/theborzet/time-tracker/internal/repository"
	"github.com/theborzet/time-tracker/internal/service"
)

func Run() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Some problems with config: %v", err)
	}

	db := database.Init(config)
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Failed to close DB: %v", err)
		}
	}()

	repo := repository.NewApiRepository(db)

	service := service.NewApiService(repo)

	handler := handler.NewApiHandler(service)

	app := fiber.New()

	routes.RegistrationRoutes(app, handler)

	if err := app.Listen(config.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
