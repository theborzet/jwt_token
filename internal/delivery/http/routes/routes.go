package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/theborzet/jwt_token/internal/delivery/http/v1/handler"
)

func RegistrationRoutes(app *fiber.App, h *handler.ApiHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	authRoutes := v1.Group("/auth")
	tokenRoutes := authRoutes.Group("/token")

	tokenRoutes.Post("/{userID}")
	tokenRoutes.Post("/refresh/{userID}")

	//Including swagger
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.json",
	}))

	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
}
