package handler

import (
	"github.com/labstack/echo/v4"
	"itemsim-server/internal/domain/item"
	"net/http"
)

type ItemHandler struct {
	itemService item.Service
}

func NewItemHandler(itemService item.Service) *ItemHandler {
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
