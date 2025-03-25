package updater

import (
	"context"
	"currency-converter/internal/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type RateFetcher interface {
	FetchRates() ([]CurrencyRate, error)
}

type Updater struct {
	Repo    *repository.CurrencyRepo
	Fetcher RateFetcher
}

func New(repo *repository.CurrencyRepo, fetcher RateFetcher) *Updater {
	return &Updater{Repo: repo, Fetcher: fetcher}
}

func (u *Updater) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				rates, err := u.Fetcher.FetchRates()
				if err != nil {
					log.Error().Err(err).Msg("Не удалось получить курсы")
					continue
				}

				for _, rate := range rates {
					err := u.Repo.UpsertExchangeRate(rate.Code, rate.Rate, rate.IsAvailable)
					if err != nil {
						log.Error().
							Err(err).
							Str("code", rate.Code).
							Msg("Ошибка при обновлении валюты")
					}
				}

				log.Info().Msg("✅ Курсы успешно обновлены")
			case <-ctx.Done():
				ticker.Stop()
				log.Info().Msg("Обновление курсов остановлено")
				return
			}
		}
	}()
}
