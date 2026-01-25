package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/programmer8760/tasks-service-api/db"
	"github.com/programmer8760/tasks-service-api/handlers/auth"
	"github.com/programmer8760/tasks-service-api/handlers/tasks"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/redis/go-redis/v9"
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

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("Redis error: ", err)
	}
	log.Println("Redis connected: ", err == nil)

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
	app.Delete("/tasks/:taskID/delete", tasks.DeleteTask(database))

	log.Fatal(app.Listen(":3000"))
}
