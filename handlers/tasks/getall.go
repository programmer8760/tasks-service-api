package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/programmer8760/tasks-service-api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GetAllTasksRequest struct {
	Token string `json:"token"`
}

func GetAllTasks(db *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		var request GetAllTasksRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.ErrBadRequest
		}

		tokenClaims, err := utils.CheckJWT(request.Token)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		userID, err := utils.GetUserIDFromJWTClaims(tokenClaims)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		var tasks []models.Task

		val, err := rdb.Get(
			ctx,
			fmt.Sprintf("tasks:user:%s", userID),
		).Result()

		if err == nil {
			json.Unmarshal([]byte(val), &tasks)
		} else {
			if err != redis.Nil {
				log.Println("Redis GET error:", err)
			}

			db.Find(
				&tasks,
				models.Task{UserID: userID},
			)

			tasksJSON, _ := json.Marshal(tasks)
			rdb.Set(
				ctx,
				fmt.Sprintf("tasks:user:%s", userID),
				tasksJSON,
				time.Minute*1,
			)
		}

		return c.JSON(fiber.Map{
			"status": 200,
			"tasks":  tasks,
		})
	}
}
