package auth

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/programmer8760/tasks-service-api/utils"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Register(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request RegisterRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.ErrBadRequest
		}

		hash, err := utils.HashPassword(request.Password)
		if err != nil {
			log.Println("bcrypt error: ", err)
			return fiber.ErrInternalServerError
		}

		user := models.User{Login: request.Login, Password: hash}

		db.Create(&user)

		return c.SendString(fmt.Sprintf("200: created new user with login \"%s\"", request.Login))
	}
}
