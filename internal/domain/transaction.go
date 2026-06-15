package domain

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	AccountID       string    `json:"account_id"`
	Amount          float64   `json:"amount"`
	MerchantID      string    `json:"merchant_id"`
	MerchantCountry string    `json:"merchant_country"`
	CreatedAt       time.Time `json:"created_at"`
}
