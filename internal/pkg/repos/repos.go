package repos

import (
	"github.com/shopspring/decimal"
	"time"
)

type AggregatedData struct {
	Day              string          `json:"day"`
	ProjectID        int             `json:"project_id"`
	NoOfTransactions int             `json:"number_of_transactions"`
	TotalUSDVolume   decimal.Decimal `json:"total_volume_usd"`
}

type Transaction struct {
	Timestamp         time.Time
	Event             string
	ProjectID         int
	CurrencySymbol    string
	CurrencyValue     decimal.Decimal
	USDConversionRate decimal.Decimal
	USDValue          decimal.Decimal
}

type MarketplaceData interface {
	InsertTransactions(transactions []Transaction) error
	GetAggregatedData() ([]AggregatedData, error)
	TruncateAllData() error
}
