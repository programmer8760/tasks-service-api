package auth

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/programmer8760/tasks-service-api/models"
	"github.com/programmer8760/tasks-service-api/utils"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request LoginRequest

		err := c.BodyParser(&request)
		if err != nil {
			return fiber.ErrBadRequest
		}

		var user models.User
		result := db.
			Where("login = ?", request.Login).
			First(&user)

		if result.Error != nil {
			log.Println("error getting data from db: ", result.Error)
			return fiber.ErrInternalServerError
		}
		if result.RowsAffected == 0 {
			return fiber.ErrNotFound
		}

		err = utils.CheckPassword(user.Password, request.Password)
		if err != nil {
			return fiber.ErrForbidden
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
		})
		secret := []byte(os.Getenv("JWT_SECRET"))
		tokenString, err := token.SignedString(secret)

		return c.JSON(fiber.Map{
			"status": 200,
			"data":   tokenString,
		})
	}
}
