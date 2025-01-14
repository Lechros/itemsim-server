package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"itemsim-server/internal/repository"
	"net/http"
	"strings"
)

func UseGearRoutes(group *echo.Group) {
	group.GET("/search", search)
	group.GET("/:id", getGearById)
	group.GET("/:id/icon", getGearIconById)
	group.GET("/:id/icon/origin", getGearIconOriginById)
	group.GET("/:id/raw-icon", getGearRawIconById)
	group.GET("/:id/raw-icon/origin", getGearRawIconOriginById)
}

func search(c echo.Context) error {
	query := c.QueryParam("query")
	query = strings.TrimSpace(query)
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query is required")
	}
	res := repository.SearchGearByName(query, 100)
	return c.JSON(http.StatusOK, res)
}

func getGearById(c echo.Context) error {
	id := c.Param("id")
	gear, ok := repository.GetGearById(id)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, gear)
}

func getGearIconById(c echo.Context) error {
	id := c.Param("id")
	url := fmt.Sprintf("https://image.itemsim.com/gears/icon/%d.png", id)
	return c.Redirect(http.StatusPermanentRedirect, url)
}

func getGearIconOriginById(c echo.Context) error {
	id := c.Param("id")
	origin, ok := repository.GetGearIconOriginById(id)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, origin)
}

func getGearRawIconById(c echo.Context) error {
	id := c.Param("id")
	url := fmt.Sprintf("https://image.itemsim.com/gears/iconRaw/%d.png", id)
	return c.Redirect(http.StatusPermanentRedirect, url)
}

func getGearRawIconOriginById(c echo.Context) error {
	id := c.Param("id")
	origin, ok := repository.GetGearRawIconOriginById(id)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, origin)
}
