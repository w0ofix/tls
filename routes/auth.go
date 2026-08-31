package routes

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/w0ofix/tls/models"
	"github.com/w0ofix/tls/utils"
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
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	var user models.User
	if err := h.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid credentials"})
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, c.IP(), c.UserAgent())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not login"})
	}

	if !slices.Contains(user.Ips.LoginIps, c.IP()) {
		user.Ips.LoginIps = append(user.Ips.LoginIps, c.IP())
		h.DB.Save(&user)
	}

	return c.JSON(fiber.Map{"type": "Bearer", "access_token": token})
}

func (h *AuthHandler) register(c fiber.Ctx) error {
	if utils.VPNCheck(c.IP()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "This action is not allowed from a suspicious network provider (VPN, Proxy, Hosting and Mobile data networks)",
		})
	}

	var body struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if len(body.Username) < 3 || len(body.Username) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Username must be between 3 and 20 characters"})
	}
	if len(body.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Password must be at least 8 characters"})
	}

	if !utils.IsValidEmail(body.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid email"})
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
	return c.SendStatus(fiber.StatusNotImplemented)
}
