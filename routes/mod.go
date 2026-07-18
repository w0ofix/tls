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
	mod.Post("/:id/report", h.reportMod)
	mod.Get("/:id/download", h.downloadMod)
	mod.Post("/report/category", h.createReportCategory)
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

	return c.JSON(fiber.Map{"success": true, "data": mods})
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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Mod not found"})
	}

	return c.JSON(fiber.Map{"success": true, "data": mod})
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
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "You are not the author of this mod"})
		}

		if err := h.DB.Delete(&mod).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Failed to delete the mod"})
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": true, "message": "Mod has been delete"})
}

func (h *ModHandler) reportMod(c fiber.Ctx) error {
	var body struct {
		Message  string `json:"message"`
		Category string `json:"category"`
	}

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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Mod not found"})
	}

	if mod.Author == claims.UserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "You cannot report your own mod"})
	}

	if err := h.DB.Where("id = ?", body.Category).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Category doesn't exist"})
	}

	report := models.ModReport{
		ModId:      id,
		ReportedBy: claims.UserID,
		Message:    body.Message,
		Category:   body.Category,
	}

	if err := h.DB.Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to report mod"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Mod has been reported"})
}

func (h *ModHandler) downloadMod(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *ModHandler) createReportCategory(c fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
	}

	jwt := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Missing token"})
	}

	claims, err := utils.ParseToken(jwt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Invalid token"})
	}

	var user models.User
	if err := h.DB.Select("role").Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "You must have administrator privileges to perform this operation"})
	}

	if user.Role == "user" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "You must have administrator privileges to perform this operation"})
	}

	category := models.ModReportCategories{
		Name: body.Name,
	}

	if err := h.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Failed to create mod report category"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Report category has been created"})
}
