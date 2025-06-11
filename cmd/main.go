package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"itemsim-server/internal/application/handler"
	"itemsim-server/internal/common/search"
	"itemsim-server/internal/domain/gear"
	"itemsim-server/internal/domain/item"
	"itemsim-server/internal/infrastructure/repository/inmemory"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://itemsim.com", "https://itemsim.pages.dev"},
		AllowMethods: []string{http.MethodGet, http.MethodOptions},
		MaxAge:       86400,
	}))

	// DI
	gearRepository := inmemory.NewGearRepository()
	itemRepository := inmemory.NewItemRepository()

	gearSearcher := search.NewSearcher[gear.Gear](gearRepository.Count())

	gearService := gear.NewGearService(gearRepository, gearSearcher)
	itemService := item.NewItemService(itemRepository)

	systemHandler := handler.NewSystemHandler()
	gearHandler := handler.NewGearHandler(gearService)
	itemHandler := handler.NewItemHandler(itemService)

	handler.RegisterRoutes(e, systemHandler, gearHandler, itemHandler)

	// Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":1323"); err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
