package route

import (
	"os"

	"github.com/apurva2307/go-api/controllers"
	"github.com/apurva2307/go-api/middlewares"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func BuildRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to my first Golang API..")
	})
	auth := app.Group("/auth")
	auth.Use(middlewares.AuthenticateToken)
	auth.Post("signup", controllers.SignupUser)
	auth.Post("login", controllers.LoginUser)
	employee := app.Group("/employee")
	employee.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	employee.Get("/getAll", controllers.GetAllEmployees)
	employee.Post("/", controllers.CreateEmployee)
	employee.Put("/:id", controllers.UpdateEmployee)
	employee.Delete("/:id", controllers.DeleteEmployee)
}
