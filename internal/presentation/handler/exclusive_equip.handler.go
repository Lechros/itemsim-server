package handler

import (
	"itemsim-server/internal/application"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ExclusiveEquipHandler struct {
	exclusiveEquipService application.ExclusiveEquipService
}

func NewExclusiveEquipHandler(exclusiveEquipService application.ExclusiveEquipService) *ExclusiveEquipHandler {
	return &ExclusiveEquipHandler{
		exclusiveEquipService: exclusiveEquipService,
	}
}

func (h *ExclusiveEquipHandler) GetAllDataAsJson(c echo.Context) error {
	result, err := h.exclusiveEquipService.GetAllDataAsJson()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
