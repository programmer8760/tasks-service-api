package tasks

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/programmer8760/tasks-service-api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CreateTaskRequest struct {
	Token       string `json:"token"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateTask(db *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request CreateTaskRequest

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

		task := models.Task{
			Title:       request.Title,
			Description: request.Description,
			Status:      models.TaskNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserID:      userID,
		}

		db.Create(&task)

		rdb.Del(c.Context(), fmt.Sprintf("tasks:user:%s", userID))

		return c.JSON(fiber.Map{
			"status": 200,
			"task":   task,
		})
	}
}
