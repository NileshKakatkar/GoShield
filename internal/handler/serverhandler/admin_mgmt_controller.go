package serverhandler

import (
	"goshield/internal/db"
	"goshield/internal/log"
	"goshield/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AdminMgmt interface {
	Authenticate(c *fiber.Ctx) error
}

type adminMgmt struct{}

func NewAdminMgmt() AdminMgmt {
	return &adminMgmt{}
}

var jwtSecret = []byte("b4c1c8f4e9d5a7b2f3c8d1e6a9b7c4d2f8e1a6c3b9d5e7f2a4c8d1b6e9f3a7c5")

func (a *adminMgmt) Authenticate(c *fiber.Ctx) error {
	var reqBody model.Admin

	if err := c.BodyParser(&reqBody); err != nil {
		log.Logger.Errorf("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	admin, err := db.NewAdminRepo().GetAdminDetails()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	isMatched := checkPassword(admin.Password, reqBody.Password)
	if !isMatched || reqBody.AdminId != admin.AdminId {
		log.Logger.Warnf("Password does not matched for %s", reqBody.AdminId)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized request",
		})
	}

	jwtToken, _ := generateJWT(admin.AdminId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   jwtToken,
	})
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	return err == nil
}

func generateJWT(adminId string) (string, error) {
	claims := jwt.MapClaims{
		"adminId": adminId,
		"sub":     adminId,
		"exp":     time.Now().Add(60 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
