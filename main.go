package main

import (
	"context"
	"log"
	"os"

	"github.com/apurva2307/go-api/route"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	mongoUri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	route.BuildRoutes(app, client)
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("This route does not exist.") // => 404 "Not Found"
	})
	log.Fatal(app.Listen(":3000"))
}
