package types

import (
	"time"

	// "github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

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
	Multiplier            string
	BuyingPower           decimal.Decimal
	ShortingEnabled       bool
	LongMarketValue       decimal.Decimal
	InitialMargin         decimal.Decimal
	MaintenanceMargin     decimal.Decimal
	CashWithdrawable      decimal.Decimal
	DaytradeCount         int64
	LastMaintenanceMargin decimal.Decimal
	DaytradingBuyingPower decimal.Decimal
	RegtBuyingPower       decimal.Decimal
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}
