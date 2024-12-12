package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
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
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "query is required")
	}
	// TODO
	return c.JSON(http.StatusOK, []string{fmt.Sprintf("gear data with %s in name", query)})
}

func getGearById(c echo.Context) error {
	id := c.Param("id")
	// TODO
	return c.JSON(http.StatusOK, fmt.Sprintf("gear data with id %s", id))
}

func getGearIconById(c echo.Context) error {
	id := c.Param("id")
	// TODO
	return c.Redirect(http.StatusPermanentRedirect, "gear icon url with id "+id)
}

func getGearIconOriginById(c echo.Context) error {
	id := c.Param("id")
	// TODO
	return c.JSON(http.StatusOK, fmt.Sprintf("gear icon origin with id %s", id))
}

func getGearRawIconById(c echo.Context) error {
	id := c.Param("id")
	// TODO
	return c.Redirect(http.StatusPermanentRedirect, "gear raw icon url with id "+id)
}

func getGearRawIconOriginById(c echo.Context) error {
	id := c.Param("id")
	// TODO
	return c.JSON(http.StatusOK, fmt.Sprintf("gear raw icon origin with id %s", id))
}
