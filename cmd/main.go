package main

import (
	"context"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
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
	// Prometheus
	e.Use(echoprometheus.NewMiddleware("itemsim"))

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

	gearSearcher := invindex.NewSearcher[gear.Gear](gearRepository.Count())

	gearService := application.NewGearService(gearRepository, gearSearcher)
	itemService := application.NewItemService(itemRepository)

	systemHandler := handler.NewSystemHandler()
	gearHandler := handler.NewGearHandler(gearService)
	itemHandler := handler.NewItemHandler(itemService)

	// Setup response cache
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(1<<16),
	)
	if err != nil {
		log.Fatal(err)
	}
	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(24*time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register routes
	handler.RegisterRoutes(e, systemHandler, gearHandler, itemHandler, cfg, cacheClient)

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
