package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

type UpdatePriceRequest struct {
	Material string  `json:"material"`
	NewPrice float64 `json:"new_price_usd"`
}

// UpdateMaterialPrice handles the admin endpoint to update local LME rates manually
func UpdateMaterialPrice(c *fiber.Ctx) error {
	var req UpdatePriceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := services.UpdateMaterialPrice(req.Material, req.NewPrice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Material price updated successfully",
	})
}

// GetMaterials returns all current materials and their prices
func GetMaterials(c *fiber.Ctx) error {
	catalog := services.GetAllMaterials()
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   catalog,
	})
}

// UpdateMachineConfig handles the admin endpoint to update machine rates
func UpdateMachineConfig(c *fiber.Ctx) error {
	var req services.MachineProfile
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := services.UpdateMachine(req.ID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Machine config updated successfully",
	})
}

// GetMachines returns all configured machines
func GetMachines(c *fiber.Ctx) error {
	catalog := services.GetAllMachines()
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   catalog,
	})
}
