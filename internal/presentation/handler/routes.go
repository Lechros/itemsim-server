package handler

import (
	"itemsim-server/internal/config"
	"itemsim-server/internal/presentation/middleware"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	cache "github.com/victorspringer/http-cache"
)

func RegisterRoutes(e *echo.Echo,
	systemHandler *SystemHandler,
	gearHandler *GearHandler,
	itemHandler *ItemHandler,
	setItemHandler *SetItemHandler,
	exclusiveEquipHandler *ExclusiveEquipHandler,
	soulHandler *SoulHandler,
	cfg *config.Config,
	cacheClient *cache.Client,
) {
	registerSystemRoutes(e, systemHandler, cfg)
	registerGearRoutes(e.Group("/gears"), gearHandler, cacheClient)
	registerItemRoutes(e.Group("/items"), itemHandler)
	registerSetItemRoutes(e.Group("/set-items"), setItemHandler, cacheClient)
	registerExclusiveEquipRoutes(e.Group("/exclusive-equips"), exclusiveEquipHandler, cacheClient)
	registerSoulRoutes(e.Group("/souls"), soulHandler, cacheClient)
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
	group.GET("/icon/origins", h.GetAllIconOrigins)
}

func registerItemRoutes(group *echo.Group, h *ItemHandler) {
	group.GET("/:id/raw-icon/origin", h.GetIconRawOrigin)
	group.GET("/raw-icon/origins", h.GetAllIconRawOrigins)
}

func registerSetItemRoutes(group *echo.Group, h *SetItemHandler, cacheClient *cache.Client) {
	group.GET("", h.GetAllDataAsJson, echo.WrapMiddleware(cacheClient.Middleware))
}

func registerExclusiveEquipRoutes(group *echo.Group, h *ExclusiveEquipHandler, cacheClient *cache.Client) {
	group.GET("", h.GetAllDataAsJson, echo.WrapMiddleware(cacheClient.Middleware))
}

func registerSoulRoutes(group *echo.Group, h *SoulHandler, cacheClient *cache.Client) {
	group.GET("", h.GetAllDataAsJson, echo.WrapMiddleware(cacheClient.Middleware))
}
