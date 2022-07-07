package route

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func BuildRoutes(app *fiber.App, client *mongo.Client) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hi air")
	})
	app.Get("/employees", func(c *fiber.Ctx) error {
		return c.SendString("hi air")
	})
	app.Post("/employee", func(c *fiber.Ctx) error {
		return c.SendString("hi air")
	})
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		return c.SendString("hi air")
	})
	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		return c.SendString("hi air")
	})
}
