package handler

import (
	"itemsim-server/internal/application"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SetItemHandler struct {
	setItemService application.SetItemService
}

func NewSetItemHandler(setItemService application.SetItemService) *SetItemHandler {
	return &SetItemHandler{
		setItemService: setItemService,
	}
}

func (h *SetItemHandler) GetAllDataAsJson(c echo.Context) error {
	result, err := h.setItemService.GetAllDataAsJson()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
