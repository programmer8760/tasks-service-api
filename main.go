package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/programmer8760/tasks-service-api/db"
	"github.com/programmer8760/tasks-service-api/types"
)

func main() {
	godotenv.Load()

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected: ", database != nil)

	database.AutoMigrate(
		&types.User{},
		&types.Task{},
		&types.Event{},
	)

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("200")
	})

	log.Fatal(app.Listen(":3000"))
}
