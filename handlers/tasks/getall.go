package tasks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/programmer8760/tasks-service-api/utils"
	"gorm.io/gorm"
)

type GetAllTasksRequest struct {
	Token string `json:"token"`
}

func GetAllTasks(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		db.Find(
			&tasks,
			models.Task{UserID: userID},
		)

		return c.JSON(fiber.Map{
			"status": 200,
			"tasks":  tasks,
		})
	}
}
