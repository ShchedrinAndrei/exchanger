package repository

import (
	"currency-converter/internal/model"
	"database/sql"
	"fmt"
)

type CurrencyRepo struct {
	db *sql.DB
}

func NewCurrencyRepo(db *sql.DB) *CurrencyRepo {
	return &CurrencyRepo{db: db}
}

func (r *CurrencyRepo) GetCurrency(code string) (*model.Currency, error) {
	var c model.Currency
	err := r.db.QueryRow(`
		SELECT code, is_available, rate, updated_at
		FROM currencies
		WHERE code = $1
	`, code).Scan(&c.Code, &c.IsAvailable, &c.Rate, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return &c, fmt.Errorf("валюта %s не найдена", code)
	}

	return &c, err
}

func (r *CurrencyRepo) UpsertExchangeRate(code string, rate float64, isAvailable bool) error {
	_, err := r.db.Exec(`
		INSERT INTO currencies (code, rate, is_available, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (code)
		DO UPDATE SET rate = EXCLUDED.rate, updated_at = NOW()
	`, code, rate, isAvailable)
	return err
}
