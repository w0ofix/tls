package routes

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"github.com/w0ofix/tls/utils"
)

type NotificationHandler struct {
	DB *gorm.DB
}

func RegisterNotificationRoutes(router fiber.Router, db *gorm.DB) {
	h := &NotificationHandler{DB: db}

	notification := router.Group("/notifications")

	notification.Get("/", h.getNotifications)
	notification.Get("/:id", h.getNotification)
	notification.Put("/:id/read", h.readNotification)
}

func (h *NotificationHandler) getNotifications(c fiber.Ctx) error {
	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	_, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.SendStatus(fiber.StatusNotImplemented)
}
func (h *NotificationHandler) getNotification(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *NotificationHandler) readNotification(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
