package handler

import (
	"itemsim-server/internal/application"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

func (h *ItemHandler) GetAllIconRawOrigins(c echo.Context) error {
	ids := c.QueryParam("id")
	if ids == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ids is required")
	}
	idList := strings.Split(ids, ",")
	result, err := h.itemService.GetAllIconRawOriginsById(idList)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "not found")
	}
	return c.JSON(http.StatusOK, result)
}
