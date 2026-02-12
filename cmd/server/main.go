package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/barkhayot/backend-test-golang-temp/pkg/apibalance/apibalancecharge"
	"github.com/barkhayot/backend-test-golang-temp/pkg/apibalance/apibalancehistory"
	"github.com/barkhayot/backend-test-golang-temp/pkg/apiskinport"
	"github.com/barkhayot/backend-test-golang-temp/pkg/appctx"
	"github.com/barkhayot/backend-test-golang-temp/pkg/balance"
	"github.com/barkhayot/backend-test-golang-temp/pkg/pgdb"
	"github.com/barkhayot/backend-test-golang-temp/pkg/skinport"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
)

const (
	appName = "backend-test"
)

func main() {
	// load env params
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	params, err := load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// infra setup
	app := fiber.New(fiber.Config{
		AppName: appName,
	})
	app.Use(logger.New())

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	db, err := pgdb.New(ctx, params.Pgdb)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pgdb.Close(db)

	if params.RunMigrations {
		if err := pgdb.RunMigrations(
			ctx,
			db,
			[]string{pgdb.TableUser, pgdb.TableUserBalanceHistory},
		); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
		log.Println("database migrations applied successfully")
	}

	// set skinport service
	skinport := skinport.NewService(
		skinport.WithFetchInterval(params.FetchInterval),
		skinport.WithFetchOnLaunch(params.FetchOnLaunch),
	)
	go skinport.Run(ctx)

	// set balance service
	balance, err := balance.NewService(db)
	if err != nil {
		log.Fatalf("failed to create balance service: %v", err)
	}

	// set services into app context
	appCtx := appctx.NewAppCtx(
		appctx.WithFiberApp(app),
		appctx.WithSkinportService(skinport),
		appctx.WithBalanceService(balance),
	)

	if err := registerRoutes(appCtx); err != nil {
		log.Fatalf("failed to register routes: %v", err)
	}

	// Start server
	go func() {
		log.Printf("HTTP server listening on %d", params.AppPort)
		if err := app.Listen(fmt.Sprintf(":%d", params.AppPort)); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("shutting down gracefully...")

	// Create new context to manage shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func registerRoutes(appCtx *appctx.AppCtx) error {
	if appCtx == nil {
		return fmt.Errorf("app context is nil")
	}

	app := appCtx.GetFiberApp()
	app.Get("/items", apiskinport.Handler(appCtx))
	app.Get("/balance/:id/history", apibalancehistory.Handler(appCtx))
	app.Post("/balance/:id/charge", apibalancecharge.Handler(appCtx))

	return nil
}
