package model

var AllowedCurrencies = map[string]struct{}{
	"EUR":  {},
	"USD":  {},
	"CNY":  {},
	"USDT": {},
	"USDC": {},
	"ETH":  {},
}

type ConvertRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type ConvertResponse struct {
	From            string  `json:"from"`
	To              string  `json:"to"`
	OriginalAmount  float64 `json:"originalAmount"`
	ConvertedAmount float64 `json:"convertedAmount"`
	Rate            float64 `json:"rate"`
}

type Currency struct {
	Code        string  `json:"code"`
	Rate        float64 `json:"rate"`
	IsAvailable bool    `json:"isAvailable"`
	UpdatedAt   string  `json:"updatedAt"`
}
