package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"itemsim-server/internal/application"
	"itemsim-server/internal/common/search/invindex"
	"itemsim-server/internal/config"
	"itemsim-server/internal/domain/gear"
	"itemsim-server/internal/infrastructure/repository/inmemory"
	"itemsim-server/internal/presentation/handler"
	"log"
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
		AllowOrigins: []string{"https://itemsim.com", "https://itemsim.pages.dev", "https://next.itemsim.com"},
		AllowMethods: []string{http.MethodGet, http.MethodOptions},
		MaxAge:       86400,
	}))

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize repositories
	gearRepository, err := inmemory.NewGearRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize gear repository: %v", err)
	}

	itemRepository, err := inmemory.NewItemRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize item repository: %v", err)
	}

	//gearSearcher := search.NewRegexSearcher[gear.Gear](gearRepository.Count())
	gearSearcher := invindex.NewSearcher[gear.Gear](gearRepository.Count())

	gearService := application.NewGearService(gearRepository, gearSearcher)
	itemService := application.NewItemService(itemRepository)

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
