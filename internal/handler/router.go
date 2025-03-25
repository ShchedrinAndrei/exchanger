package handler

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Get("/convert", h.ConvertHandler)
}
