package main

import (
	"context"
	"log"
	"os"

	"github.com/apurva2307/go-api/db"
	"github.com/apurva2307/go-api/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	mongoUri := os.Getenv("MONGO_URI")
	mg, err := db.ConnectDb(mongoUri, "go-api")
	if err != nil {
		panic(err)
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://goapi.apurvasingh.dev, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	route.BuildRoutes(app)
	defer mg.Client.Disconnect(context.Background())
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("This route does not exist.") // => 404 "Not Found"
	})
	log.Fatal(app.Listen(":3000"))
}
