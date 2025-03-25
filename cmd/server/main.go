package main

import (
	"context"
	"currency-converter/internal/repository"
	"currency-converter/internal/service"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"currency-converter/config"
	"currency-converter/internal/db"
	"currency-converter/internal/handler"
	"currency-converter/internal/updater"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("DB connection failed")
	}

	if err := goose.Up(dbConn, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("DB migration failed")
	}

	app := fiber.New()

	repo := repository.NewCurrencyRepo(dbConn)
	converter := service.NewConverterService(repo)
	h := handler.NewHandler(converter)
	fetcher := updater.NewFastForexFetcher()
	updater.New(repo, fetcher).Start(context.Background(), time.Duration(cfg.FetchIntervalMinute)*time.Minute)
	handler.RegisterRoutes(app, h)

	log.Info().Str("port", cfg.ServerPort).Msg("🚀 Server is running")
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatal().Err(err).Msg("Fiber server error")
	}
}
