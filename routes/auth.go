package routes

import (
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/w0ofix/tls/models"
)

type AuthHandler struct {
	DB *gorm.DB
}

func RegisterAuthRoutes(router fiber.Router, db *gorm.DB) {
	h := &AuthHandler{DB: db}

	auth := router.Group("/auth")

	auth.Post("/login", h.login)
	auth.Post("/register", h.register)
	auth.Post("/logout", h.logout)
}

func (h *AuthHandler) login(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"token": "fake-jwt"})
}

func (h *AuthHandler) register(c fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not create user"})
	}

	user := models.User{
		Email:    body.Email,
		Username: body.Username,
		Password: string(hashedPassword),
		Ips: models.Ips{
			RegisteredIp: c.IP(),
			LoginIps:     []string{},
		},
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not create user"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "User registered successfully"})
}

func (h *AuthHandler) logout(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
