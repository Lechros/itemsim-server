package handler

import (
	"github.com/labstack/echo/v4"
	"itemsim-server/internal/domain/gear"
	"net/http"
	"strconv"
)

type GearHandler struct {
	gearService gear.Service
}

func NewGearHandler(gearService gear.Service) *GearHandler {
	return &GearHandler{
		gearService: gearService,
	}
}

func (h *GearHandler) Search(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query is required")
	}
	results := h.gearService.SearchByName(query)
	return c.JSON(http.StatusOK, results)
}

func (h *GearHandler) GetData(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	result, err := h.gearService.GetDataById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *GearHandler) GetIconOrigin(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	result, err := h.gearService.GetIconOriginById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	return c.JSON(http.StatusOK, result)
}
