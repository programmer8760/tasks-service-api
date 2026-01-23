package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/programmer8760/tasks-service-api/db"
	"github.com/programmer8760/tasks-service-api/handlers/auth"
	"github.com/programmer8760/tasks-service-api/handlers/tasks"
	"github.com/programmer8760/tasks-service-api/models"
)

func main() {
	godotenv.Load()

	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected: ", database != nil)

	database.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Event{},
	)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("200")
	})

	app.Post("/register", auth.Register(database))
	app.Post("/login", auth.Login(database))

	app.Post("/tasks/create", tasks.CreateTask(database))
	app.Get("/tasks", tasks.GetAllTasks(database))
	app.Get("/tasks/:taskID", tasks.GetOneTask(database))
	app.Put("/tasks/:taskID/update", tasks.UpdateTask(database))

	log.Fatal(app.Listen(":3000"))
}
