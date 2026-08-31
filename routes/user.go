package routes

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/w0ofix/tls/models"
	"github.com/w0ofix/tls/utils"
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

	_, err := utils.ParseToken(jwt, c.IP(), c.UserAgent())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

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

	claims, err := utils.ParseToken(jwt, c.IP(), c.UserAgent())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": err.Error()})
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
	var body struct {
		Email    *string `json:"email"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Bio      *string `json:"bio"`
	}

	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	claims, err := utils.ParseToken(jwt, c.IP(), c.UserAgent())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	id := c.Params("id")
	if id != "me" {
		return c.SendStatus(fiber.StatusNotImplemented)
	}

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	updates := map[string]interface{}{}
	if body.Email != nil {
		if !utils.IsValidEmail(*body.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid email"})
		}

		updates["email"] = *body.Email
	}
	if body.Username != nil {
		if len(*body.Username) < 3 || len(*body.Username) > 20 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Username must be between 3 and 20 characters"})
		}

		updates["username"] = *body.Username
	}
	if body.Bio != nil {
		if len(*body.Bio) > 240 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Bio must be less than 240 characters"})
		}

		updates["bio"] = *body.Bio
	}
	if body.Password != nil {
		if len(*body.Password) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Password must be at least 8 characters"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not update password"})
		}
		updates["password"] = string(hashedPassword)
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "No fields to update"})
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", claims.UserID).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not update user"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "User updated successfully"})
}
