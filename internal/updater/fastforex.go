package updater

import (
	"currency-converter/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type CurrencyRate struct {
	Code        string
	Rate        float64
	IsAvailable bool
}

type FastForexFetcher struct {
	client *http.Client
	apiKey string
}

func NewFastForexFetcher() *FastForexFetcher {
	return &FastForexFetcher{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: os.Getenv("FASTFOREX_API_KEY"),
	}
}

func (f *FastForexFetcher) FetchRates() ([]CurrencyRate, error) {
	url := fmt.Sprintf("%s?api_key=%s", os.Getenv("FASTFOREX_URL"), f.apiKey)

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к fastforex.io: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fastforex вернул статус %d", resp.StatusCode)
	}

	var body struct {
		Base    string             `json:"base"`
		Results map[string]float64 `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ fastforex: %w", err)
	}

	var result []CurrencyRate
	found := make(map[string]bool)

	for code, rate := range body.Results {
		upper := strings.ToUpper(code)
		if _, ok := model.AllowedCurrencies[upper]; ok {
			result = append(result, CurrencyRate{
				Code:        upper,
				Rate:        rate,
				IsAvailable: true,
			})
			found[upper] = true
		}
	}

	for code := range model.AllowedCurrencies {
		if !found[code] {
			result = append(result, CurrencyRate{
				Code:        code,
				Rate:        0,
				IsAvailable: false,
			})
		}
	}

	return result, nil
}
