package handler

import (
	"itemsim-server/internal/application"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SoulHandler struct {
	soulService application.SoulService
}

func NewSoulHandler(soulService application.SoulService) *SoulHandler {
	return &SoulHandler{
		soulService: soulService,
	}
}

func (h *SoulHandler) GetAllDataAsJson(c echo.Context) error {
	result, err := h.soulService.GetAllDataAsJson()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (h *SoulHandler) Search(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query is required")
	}
	results, err := h.soulService.SearchByName(query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, results)
}

func (h *SoulHandler) GetData(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	result, err := h.soulService.GetDataById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	return c.JSON(http.StatusOK, result)
}
