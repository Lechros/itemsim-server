package handler

import (
	"itemsim-server/internal/application"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type GearHandler struct {
	gearService application.GearService
}

func NewGearHandler(gearService application.GearService) *GearHandler {
	return &GearHandler{
		gearService: gearService,
	}
}

func (h *GearHandler) Search(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query is required")
	}
	prefix := c.QueryParam("type")
	var prefixInt *int = nil
	if prefix != "" {
		_prefixInt, err := strconv.Atoi(prefix)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid type")
		}
		prefixInt = &_prefixInt
	}
	results, err := h.gearService.SearchByName(query, prefixInt)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
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

func (h *GearHandler) GetAllData(c echo.Context) error {
	ids := c.QueryParam("id")
	if ids == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ids is required")
	}
	idsStrList := strings.Split(ids, ",")
	idList := make([]int, 0, len(idsStrList))
	for _, idStr := range idsStrList {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid ids")
		}
		idList = append(idList, id)
	}
	results, err := h.gearService.GetAllDataById(idList)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, results)
}

func (h *GearHandler) GetHash(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	result, err := h.gearService.GetHashById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	return c.JSON(http.StatusOK, result)
}

func (h *GearHandler) GetAllHashes(c echo.Context) error {
	ids := c.QueryParam("id")
	if ids == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ids is required")
	}
	idsStrList := strings.Split(ids, ",")
	idList := make([]int, 0, len(idsStrList))
	for _, idStr := range idsStrList {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid ids")
		}
		idList = append(idList, id)
	}
	results, err := h.gearService.GetAllHashesById(idList)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, results)
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

func (h *GearHandler) GetAllIconOrigins(c echo.Context) error {
	ids := c.QueryParam("id")
	if ids == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ids is required")
	}
	idsStrList := strings.Split(ids, ",")
	idList := make([]int, 0, len(idsStrList))
	for _, idStr := range idsStrList {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid ids")
		}
		idList = append(idList, id)
	}
	result, err := h.gearService.GetAllIconOriginsById(idList)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "not found")
	}
	return c.JSON(http.StatusOK, result)
}
