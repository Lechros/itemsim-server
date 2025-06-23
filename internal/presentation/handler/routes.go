package handler

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	cache "github.com/victorspringer/http-cache"
	"itemsim-server/internal/config"
	"itemsim-server/internal/presentation/middleware"
)

func RegisterRoutes(e *echo.Echo,
	systemHandler *SystemHandler,
	gearHandler *GearHandler,
	itemHandler *ItemHandler,
	cfg *config.Config,
	cacheClient *cache.Client,
) {
	registerSystemRoutes(e, systemHandler, cfg)
	registerGearRoutes(e.Group("/gears"), gearHandler, cacheClient)
	registerItemRoutes(e.Group("/items"), itemHandler)
}

func registerSystemRoutes(e *echo.Echo, h *SystemHandler, cfg *config.Config) {
	e.GET("/health", h.Healthcheck)
	e.GET("/metrics", echoprometheus.NewHandler(), middleware.BearerAuth(cfg.MetricsPassword))
}

func registerGearRoutes(group *echo.Group, h *GearHandler, cacheClient *cache.Client) {
	group.GET("/search", h.Search, echo.WrapMiddleware(cacheClient.Middleware))
	group.GET("", h.GetAllData)
	group.GET("/:id", h.GetData)
	group.GET("/:id/icon/origin", h.GetIconOrigin)
}

func registerItemRoutes(group *echo.Group, h *ItemHandler) {
	group.GET("/:id/raw-icon/origin", h.GetIconRawOrigin)
}
