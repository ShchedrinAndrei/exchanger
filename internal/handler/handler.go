package handler

import (
	"currency-converter/internal/model"
	"currency-converter/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type Handler struct {
	Converter *service.ConverterService
}

func NewHandler(converter *service.ConverterService) *Handler {
	return &Handler{Converter: converter}
}

func (h *Handler) ConvertHandler(c *fiber.Ctx) error {
	from := strings.ToUpper(c.Query("from"))
	to := strings.ToUpper(c.Query("to"))
	amountStr := c.Query("amount")

	var errs []model.FieldError

	if from == "" {
		errs = append(errs, model.FieldError{Field: "from", Message: "Обязательное поле"})
	}
	if to == "" {
		errs = append(errs, model.FieldError{Field: "to", Message: "Обязательное поле"})
	}

	if _, ok := model.AllowedCurrencies[from]; !ok {
		errs = append(errs, model.FieldError{Field: "from", Message: "Валюта не поддерживается"})
	}
	if _, ok := model.AllowedCurrencies[to]; !ok {
		errs = append(errs, model.FieldError{Field: "to", Message: "Валюта не поддерживается"})
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if amountStr == "" {
		errs = append(errs, model.FieldError{Field: "amount", Message: "Обязательное поле"})
	} else if err != nil || amount <= 0 {
		errs = append(errs, model.FieldError{Field: "amount", Message: "Должно быть числом больше 0"})
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Errors: errs,
		})
	}

	req := model.ConvertRequest{
		From:   from,
		To:     to,
		Amount: amount,
	}

	result, err := h.Converter.Convert(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Errors: []model.FieldError{
				{Field: "general", Message: err.Error()},
			},
		})
	}

	return c.JSON(result)
}
