package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/theborzet/time-tracker/internal/delivery/http/handler"
)

func RegistrationRoutes(app *fiber.App, h *handler.ApiHandler) {
	userRoutes := app.Group("/user")
	taskRoutes := app.Group("/task")

	userRoutes.Get("/tasks", h.GetUserTasks)
	userRoutes.Get("/", h.GetUsers)
	userRoutes.Post("/create", h.CreateUser)
	userRoutes.Put("/update", h.UpdateUser)
	userRoutes.Delete("/:id", h.DeleteUser)

	taskRoutes.Post("/start", h.StartTask)
	taskRoutes.Post("/end", h.EndTask)

	//Including swagger
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.json",
	}))

	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
}
