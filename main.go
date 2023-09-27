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
)

type Account struct {
	gorm.Model
	ID                    uint
	AccountID             string
	AccountBlocked        bool
	AccountNumber         string
	crypto_status         string
	currency              string
	cash                  decimal.Decimal
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

func getAccount(client alpaca.Client, db gorm.DB) {
	// Get account information
	acct, err := client.GetAccount()

	var account Account = Account{
		AccountID:             acct.ID,
		AccountBlocked:        acct.AccountBlocked,
		AccountNumber:         acct.AccountNumber,
		PortfolioValue:        acct.PortfolioValue,
		crypto_status:         acct.Status,
		currency:              acct.Currency,
		cash:                  acct.Cash,
		PatternDayTrader:      acct.PatternDayTrader,
		TradingBlocked:        acct.TradingBlocked,
		TransfersBlocked:      acct.TransfersBlocked,
		ShortMarketValue:      acct.ShortMarketValue,
		Equity:                acct.Equity,
		LastEquity:            acct.LastEquity,
		Multiplier:            acct.Multiplier,
		BuyingPower:           acct.BuyingPower,
		ShortingEnabled:       acct.ShortingEnabled,
		LongMarketValue:       acct.LongMarketValue,
		InitialMargin:         acct.InitialMargin,
		MaintenanceMargin:     acct.MaintenanceMargin,
		CashWithdrawable:      acct.CashWithdrawable,
		DaytradeCount:         acct.DaytradeCount,
		LastMaintenanceMargin: acct.LastMaintenanceMargin,
		DaytradingBuyingPower: acct.DaytradingBuyingPower,
		RegtBuyingPower:       acct.RegTBuyingPower,
		CreatedAt:             acct.CreatedAt,
		UpdatedAt:             time.Now(),
	}

	if err != nil {
		// Print error
		fmt.Printf("Error getting account information: %v", err)
	} else {
		// Print account information
		fmt.Printf("Account ID: %+v\n", acct.ID)
		fmt.Printf("Account: %+v\n", *acct)

		// FirstOrCreate
		result := db.Where(Account{AccountID: account.AccountID}).FirstOrCreate(&account)

		fmt.Printf("Query Result: %+v\n", result)
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

	dsn := os.Getenv("dbuser") + ":" + os.Getenv("dbpass") + "@tcp(" + os.Getenv("dbhost") + ":" + os.Getenv("dbport") + ")/" + os.Getenv("dbname") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

type Order struct {
	gorm.Model
	OrderId        string
	ClientOrderId  string
	CreatedAt      string
	UpdatedAt      string
	SubmittedAt    string
	FilledAt       string
	ExpiredAt      string
	CanceledAt     string
	FailedAt       string
	ReplacedAt     string
	ReplacedBy     string
	Replaces       string
	AssetId        string
	Symbol         string
	AssetClass     string
	Notional       string
	Qty            string
	FilledQty      string
	FilledAvgPrice string
	OrderType      string
	Type           string
	Side           string
	TimeInForce    string
	LimitPrice     string
	Status         string
}

func placeOrder(client alpaca.Client, db gorm.DB) (alpaca.Order, error) {
	symbol := "BTC/USD"
	qty := decimal.NewFromInt(1)
	side := alpaca.Side("buy")
	orderType := alpaca.OrderType("market")
	timeInForce := alpaca.TimeInForce("day")

	// Placing an order with the parameters set previously
	order, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
		AssetKey:    &symbol,
		Qty:         &qty,
		Side:        side,
		Type:        orderType,
		TimeInForce: timeInForce,
	})
	if err != nil {
		// Print error
		fmt.Printf("Failed to place order: %v\n", err)
	} else {
		// Print resulting order object
		fmt.Printf("Order succesfully sent:\n%+v\n", *order)
	}

	return *order, err
}
