package routes

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"github.com/w0ofix/tls/models"
	"github.com/w0ofix/tls/utils"
)

type ModHandler struct {
	DB *gorm.DB
}

func RegisterModRoutes(router fiber.Router, db *gorm.DB) {
	h := &ModHandler{DB: db}

	mod := router.Group("/mods")

	mod.Get("/", h.getMods)
	mod.Post("/", h.postMod)
	mod.Get("/:id", h.getMod)
	mod.Delete("/:id", h.deleteMod)
}

func (h *ModHandler) getMods(c fiber.Ctx) error {
	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	_, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	var mods []models.Mod
	if err := h.DB.Select("id", "name", "description", "author", "game", "category", "price", "image", "path").Where("status = 1").Find(&mods).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not fetch mods"})
	}

	return c.JSON(fiber.Map{"status": true, "data": mods})
}

func (h *ModHandler) getMod(c fiber.Ctx) error {
	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	_, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	id := c.Params("id")

	var mod models.Mod
	if err := h.DB.Select("id", "name", "description", "author", "game", "category", "price", "image", "path").Where("id = ? AND status = 1", id).First(&mod).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Mod not found"})
	}

	return c.JSON(fiber.Map{"status": true, "data": mod})
}

func (h *ModHandler) postMod(c fiber.Ctx) error {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Game        string `json:"game"`
		Category    string `json:"category"`
		Price       int    `json:"price"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	claims, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	mod := models.Mod{
		Name:        body.Name,
		Description: body.Description,
		Author:      claims.UserID,
		Game:        body.Game,
		Category:    body.Category,
		Price:       body.Price,
	}

	if err := h.DB.Create(&mod).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Could not post mod"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Mod has been posted"})
}

func (h *ModHandler) deleteMod(c fiber.Ctx) error {
	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	claims, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	id := c.Params("id")

	var mod models.Mod
	if err := h.DB.Where("id = ?", id).First(&mod); err != nil {
		if mod.Author != claims.UserID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": false, "message": "You are not the author of this mod"})
		}

		if err := h.DB.Delete(&mod).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": false, "message": "Failed to delete the mod"})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": true, "message": "Mod has been delete"})
}
