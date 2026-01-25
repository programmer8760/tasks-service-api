package tasks

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/programmer8760/tasks-service-api/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UpdateTaskRequest struct {
	Token       string        `json:"token"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      models.Status `json:"status"`
}

func UpdateTask(db *gorm.DB, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request UpdateTaskRequest

		taskID := c.Params("taskID")

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

		if request.Title == "" && request.Description == "" && request.Status == "" {
			return fiber.ErrBadRequest
		}
		if !request.Status.IsValid() {
			return fiber.ErrBadRequest
		}

		var task models.Task

		result := db.
			Where("id = ? AND user_id = ?", taskID, userID).
			Find(&task)

		if result.Error != nil {
			log.Println("error getting data from db: ", result.Error)
			return fiber.ErrInternalServerError
		}
		if result.RowsAffected == 0 {
			return fiber.ErrNotFound
		}

		if request.Description != "" {
			task.Description = request.Description
		}
		if request.Title != "" {
			task.Title = request.Title
		}
		if request.Status != "" {
			task.Status = request.Status
		}
		task.UpdatedAt = time.Now()

		db.Save(&task)

		go func(userID uint) {
			rdb.Del(c.Context(), fmt.Sprintf("tasks:user:%s", userID))
		}(userID)

		return c.JSON(fiber.Map{
			"status": 200,
			"task":   task,
		})
	}
}
