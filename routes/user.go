package routes

import (
	"time"
	"strings"

	"gorm.io/gorm"
	"github.com/gofiber/fiber/v3"

	"github.com/w0ofix/tls/utils"
	"github.com/w0ofix/tls/models"
)

type UserHandler struct {
	DB *gorm.DB
}

type PublicUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func RegisterUserRoutes(router fiber.Router, db *gorm.DB) {
	h := &UserHandler{DB: db}

	user := router.Group("/users")

	user.Get("/", h.getUsers)
	user.Get("/:id", h.getUser)
	user.Put("/:id", h.updateUser)
}

func (h *UserHandler) getUsers(c fiber.Ctx) error {
	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	claims, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	_ = claims

	var users []models.User
	if err := h.DB.Select("id", "email", "username", "bio", "avatar", "role", "created_at", "updated_at").Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not fetch users"})
	}

	filtered_users := make([]PublicUser, len(users))
	for i, u := range users {
		filtered_users[i] = PublicUser{
			ID:        u.ID,
			Username:  u.Username,
			Bio:       u.Bio,
			Avatar:    u.Avatar,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		}
	}

	return c.JSON(fiber.Map{"success": true, "data": filtered_users})
}

func (h *UserHandler) getUser(c fiber.Ctx) error {
	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	claims, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	id := c.Params("id")

	if id == "me" {
		var user models.User
		if err := h.DB.Select("id", "email", "username", "bio", "avatar", "role", "created_at", "updated_at").Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "User not found"})
		}

		return c.JSON(fiber.Map{"success": true, "data": user})
	}

	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *UserHandler) updateUser(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
