package service

import (
	"currency-converter/internal/model"
)

type CurrencyRepository interface {
	GetCurrency(code string) (*model.Currency, error)
	UpsertExchangeRate(code string, rate float64, isAvailable bool) error
}

type ConverterService struct {
	repo CurrencyRepository
}

func NewConverterService(repo CurrencyRepository) *ConverterService {
	return &ConverterService{repo: repo}
}

func (s *ConverterService) Convert(req model.ConvertRequest) (*model.ConvertResponse, error) {
	curFrom, err := s.repo.GetCurrency(req.From)
	if err != nil {
		return nil, err
	}

	curTo, err := s.repo.GetCurrency(req.To)
	if err != nil {
		return nil, err
	}

	return &model.ConvertResponse{
		From:            req.From,
		To:              req.To,
		OriginalAmount:  req.Amount,
		ConvertedAmount: req.Amount * curTo.Rate / curFrom.Rate,
		Rate:            curTo.Rate / curFrom.Rate,
	}, nil
}
