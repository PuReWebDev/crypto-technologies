package types

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Account is a struct to populate the Alpaca account type.
type Account struct {
	gorm.Model
	ID                    uint
	AccountID             string
	AccountBlocked        bool
	AccountNumber         string
	Status                string
	Currency              string
	Cash                  decimal.Decimal
	PortfolioValue        decimal.Decimal
	PatternDayTrader      bool
	TradingBlocked        bool
	TransfersBlocked      bool
	ShortMarketValue      decimal.Decimal
	Equity                decimal.Decimal
	LastEquity            decimal.Decimal
	Multiplier            decimal.Decimal
	BuyingPower           decimal.Decimal
	ShortingEnabled       bool
	LongMarketValue       decimal.Decimal
	InitialMargin         decimal.Decimal
	MaintenanceMargin     decimal.Decimal
	CryptoStatus          string
	DaytradeCount         int64
	LastMaintenanceMargin decimal.Decimal
	DaytradingBuyingPower decimal.Decimal
	RegtBuyingPower       decimal.Decimal
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

type Order struct {
	gorm.Model
	OrderId        string
	ClientOrderId  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	SubmittedAt    time.Time
	FilledAt       *time.Time
	ExpiredAt      *time.Time
	CanceledAt     *time.Time
	FailedAt       *time.Time
	ReplacedAt     *time.Time
	ReplacedBy     string
	Replaces       string
	AssetId        string
	Symbol         string
	AssetClass     string
	Notional       *decimal.Decimal
	Qty            *decimal.Decimal
	FilledQty      decimal.Decimal
	FilledAvgPrice *decimal.Decimal
	OrderType      string
	Type           string
	Side           string
	TimeInForce    string
	LimitPrice     *decimal.Decimal
	Status         string
}
