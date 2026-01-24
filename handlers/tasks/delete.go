package tasks

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/programmer8760/tasks-service-api/utils"
	"gorm.io/gorm"
)

type DeleteTaskRequest struct {
	Token string `json:"token"`
}

func DeleteTask(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request DeleteTaskRequest

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

		var task models.Task
		result := db.
			Where("id = ? AND user_id = ?", taskID, userID).
			Delete(&task)

		if result.Error != nil {
			log.Println("error getting data from db: ", result.Error)
			return fiber.ErrInternalServerError
		}
		if result.RowsAffected == 0 {
			return fiber.ErrNotFound
		}

		return c.JSON(fiber.Map{
			"status": 200,
		})
	}
}
