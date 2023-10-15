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
	ReplacedBy     *string
	Replaces       *string
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

type ListOrdersRequest struct {
	Status    *string    `json:"status"`
	After     *time.Time `json:"after"`
	Until     *time.Time `json:"until"`
	Limit     *int       `json:"limit"`
	Direction *string    `json:"direction"`
	Nested    *bool      `json:"nested"`
	Symbols   *string    `json:"symbols"`
	Side      *string    `json:"side"`
}

type Side string
type TimeInForce string
type OrderClass string
type OrderType string

type TakeProfit struct {
	LimitPrice *decimal.Decimal `json:"limit_price"`
}

type StopLoss struct {
	LimitPrice *decimal.Decimal `json:"limit_price"`
	StopPrice  *decimal.Decimal `json:"stop_price"`
}

type PlaceOrderRequest struct {
	AccountID     string           `json:"-"`
	AssetKey      *string          `json:"symbol"`
	Qty           *decimal.Decimal `json:"qty"`
	Notional      *decimal.Decimal `json:"notional"`
	Side          Side             `json:"side"`
	Type          OrderType        `json:"type"`
	TimeInForce   TimeInForce      `json:"time_in_force"`
	LimitPrice    *decimal.Decimal `json:"limit_price"`
	ExtendedHours bool             `json:"extended_hours"`
	StopPrice     *decimal.Decimal `json:"stop_price"`
	ClientOrderID string           `json:"client_order_id"`
	OrderClass    OrderClass       `json:"order_class"`
	TakeProfit    *TakeProfit      `json:"take_profit"`
	StopLoss      *StopLoss        `json:"stop_loss"`
	TrailPrice    *decimal.Decimal `json:"trail_price"`
	TrailPercent  *decimal.Decimal `json:"trail_percent"`
}

type Position struct {
	AssetID        string           `json:"asset_id"`
	Symbol         string           `json:"symbol"`
	Exchange       string           `json:"exchange"`
	Class          string           `json:"asset_class"`
	AccountID      string           `json:"account_id"`
	EntryPrice     decimal.Decimal  `json:"avg_entry_price"`
	Qty            decimal.Decimal  `json:"qty"`
	Side           string           `json:"side"`
	MarketValue    *decimal.Decimal `json:"market_value"`
	CostBasis      decimal.Decimal  `json:"cost_basis"`
	UnrealizedPL   *decimal.Decimal `json:"unrealized_pl"`
	UnrealizedPLPC *decimal.Decimal `json:"unrealized_plpc"`
	CurrentPrice   *decimal.Decimal `json:"current_price"`
	LastdayPrice   *decimal.Decimal `json:"lastday_price"`
	ChangeToday    *decimal.Decimal `json:"change_today"`
}

type CryptoQuote struct {
	gorm.Model
	Symbol    string
	Exchange  string
	BidPrice  float64
	BidSize   float64
	AskPrice  float64
	AskSize   float64
	Timestamp time.Time
}

type BtcPrice struct {
	gorm.Model
	Type    string    `json:"type"`
	Time    time.Time `json:"time"`
	Product string    `json:"product_id"`
	Price   string    `json:"price"`
}
