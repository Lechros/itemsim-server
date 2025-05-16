package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"itemsim-server/internal/repository"
	"net/http"
)

func UseItemRoutes(group *echo.Group) {
	group.GET("/:id/raw-icon", gearItemRawIconById)
	group.GET("/:id/raw-icon/origin", getItemRawIconOriginById)
}

func gearItemRawIconById(c echo.Context) error {
	id := c.Param("id")
	url := fmt.Sprintf("https://image.itemsim.com/origins/iconRaw/%d.png", id)
	return c.Redirect(http.StatusPermanentRedirect, url)
}

func getItemRawIconOriginById(c echo.Context) error {
	id := c.Param("id")
	origin, ok := repository.GetItemRawIconOriginById(id)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=86400")
	return c.JSON(http.StatusOK, origin)
}
