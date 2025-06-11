package handler

import (
	"github.com/labstack/echo/v4"
	"itemsim-server/internal/application"
	"net/http"
)

type ItemHandler struct {
	itemService application.ItemService
}

func NewItemHandler(itemService application.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

func (h *ItemHandler) GetIconRawOrigin(c echo.Context) error {
	id := c.Param("id")
	result, err := h.itemService.GetIconRawOriginById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	return c.JSON(http.StatusOK, result)
}
