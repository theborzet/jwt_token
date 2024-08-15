package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	config "github.com/theborzet/jwt_token/configs"
	database "github.com/theborzet/jwt_token/internal/db"
	"github.com/theborzet/jwt_token/internal/delivery/http/routes"
	"github.com/theborzet/jwt_token/internal/delivery/http/v1/handler"
	"github.com/theborzet/jwt_token/internal/repository"
	"github.com/theborzet/jwt_token/internal/service"
	"github.com/theborzet/jwt_token/pkg/auth"
	"github.com/theborzet/jwt_token/pkg/migrator"
)

func Run() {
	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	config, err := config.LoadConfig()
	if err != nil {
		infoLogger.Fatalf("Some problems with config: %v", err)
	}

	db := database.Init(config)
	defer func() {
		if err := database.Close(db); err != nil {
			infoLogger.Printf("Failed to close DB: %v", err)
		}
	}()

	if err := migrator.RunDatabaseMigrations(db); err != nil {
		infoLogger.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewApiRepository(db, debugLogger)

	tokenMananger := auth.NewManager(config)

	service := service.NewApiService(repo, debugLogger, tokenMananger)

	handler := handler.NewApiHandler(service, debugLogger)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	routes.RegistrationRoutes(app, handler)

	if err := app.Listen(":" + config.Port); err != nil {
		infoLogger.Fatalf("Error starting server: %v", err)
	}

}
