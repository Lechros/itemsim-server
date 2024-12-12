package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UseItemRoutes(group *echo.Group) {
	group.GET("/:id/raw-icon", gearItemRawIconById)
	group.GET("/:id/raw-icon/origin", getItemRawIconOriginById)
}

func gearItemRawIconById(c echo.Context) error {
	id := c.Param("id")
	// TODO
	return c.Redirect(http.StatusPermanentRedirect, "item raw icon url with id "+id)
}

func getItemRawIconOriginById(c echo.Context) error {
	id := c.Param("id")
	// TODO
	return c.JSON(http.StatusOK, fmt.Sprintf("item raw icon origin with id %s", id))
}
