package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	ID                          uint `gorm:"primaryKey"`
	AccountID                   string
	account_blocked             bool
	account_number              string
	crypto_status               string
	currency                    string
	cash                        decimal.Decimal
	portfolio_value             decimal.Decimal
	non_marginable_buying_power string
	accrued_fees                string
	pending_transfer_in         string
	pending_transfer_out        string
	pattern_day_trader          bool
	trade_suspended_by_user     bool
	trading_blocked             bool
	transfers_blocked           bool
	short_market_value          decimal.Decimal
	equity                      decimal.Decimal
	last_equity                 decimal.Decimal
	multiplier                  string
	buying_power                decimal.Decimal
	shorting_enabled            bool
	long_market_value           decimal.Decimal
	initial_margin              decimal.Decimal
	maintenance_margin          decimal.Decimal
	sma                         string
	daytrade_count              int64
	last_maintenance_margin     decimal.Decimal
	daytrading_buying_power     decimal.Decimal
	regt_buying_power           decimal.Decimal
	created_at                  time.Time
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	DeletedAt                   gorm.DeletedAt `gorm:"index"`
}

func getAccount(client alpaca.Client, db gorm.DB) {
	// Get account information
	acct, err := client.GetAccount()

	var account Account = Account{
		AccountID:       *&acct.ID,
		account_blocked: *&acct.AccountBlocked,
		account_number:  *&acct.AccountNumber,
		portfolio_value: *&acct.PortfolioValue,
		crypto_status:   *&acct.Status,
		currency:        *&acct.Currency,
		cash:            *&acct.Cash,
		// non_marginable_buying_power: *&acct.
		// accrued_fees: *&acct.
		// pending_transfer_in: *&acct.
		pattern_day_trader: *&acct.PatternDayTrader,
		// trade_suspended_by_user: *&acct.
		trading_blocked:    *&acct.TradingBlocked,
		transfers_blocked:  *&acct.TransfersBlocked,
		short_market_value: *&acct.ShortMarketValue,
		equity:             *&acct.Equity,
		last_equity:        *&acct.LastEquity,
		multiplier:         *&acct.Multiplier,
		buying_power:       *&acct.BuyingPower,
		shorting_enabled:   *&acct.ShortingEnabled,
		long_market_value:  *&acct.LongMarketValue,
		initial_margin:     *&acct.InitialMargin,
		maintenance_margin: *&acct.MaintenanceMargin,
		// sma: *&acct.
		daytrade_count:          *&acct.DaytradeCount,
		last_maintenance_margin: *&acct.LastMaintenanceMargin,
		daytrading_buying_power: *&acct.DaytradingBuyingPower,
		regt_buying_power:       *&acct.RegTBuyingPower,
		created_at:              *&acct.CreatedAt,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	if err != nil {
		// Print error
		fmt.Printf("Error getting account information: %v", err)
	} else {
		// Print account information
		fmt.Printf("Account ID: %+v\n", *&acct.ID)
		fmt.Printf("Account: %+v\n", *acct)

		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&account)
	}
}

func main() {

	db := loadEnvironment()

	client := prepAlpaca()

	getAccount(client, *db)
}

func loadEnvironment() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Migrate the schema
	db.AutoMigrate(&Account{})

	return db
}

func prepAlpaca() alpaca.Client {
	// Instantiating new Alpaca paper trading client
	// Alternatively, you can set your API key and secret using environment
	// variables named APCA_API_KEY_ID and APCA_API_SECRET_KEY respectively
	// Remove for live trading
	apiKey := os.Getenv("apiKey")
	apiSecret := os.Getenv("apiSecret")
	baseURL := os.Getenv("baseURL")

	client := alpaca.NewClient(alpaca.ClientOpts{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		BaseURL:   baseURL,
	})

	return client
}
