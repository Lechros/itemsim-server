package handler

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo,
	systemHandler *SystemHandler,
	gearHandler *GearHandler,
	itemHandler *ItemHandler) {

	registerSystemRoutes(e, systemHandler)
	registerGearRoutes(e.Group("/gears"), gearHandler)
	registerItemRoutes(e.Group("/items"), itemHandler)
}

func registerSystemRoutes(e *echo.Echo, h *SystemHandler) {
	e.GET("/health", h.Healthcheck)
}

func registerGearRoutes(group *echo.Group, h *GearHandler) {
	group.GET("/search", h.Search)
	group.GET("/:id", h.GetData)
	group.GET("/:id/icon/origin", h.GetIconOrigin)
}

func registerItemRoutes(group *echo.Group, h *ItemHandler) {
	group.GET("/:id/raw-icon/origin", h.GetIconRawOrigin)
}
