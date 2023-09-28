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
	placeOrder(client, *db)
	listPositions(client, *db)
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
	db.AutoMigrate(&Account{}, &Order{})

	return db
}

func prepAlpaca() alpaca.Client {
	apiKey := os.Getenv("APCA_API_KEY_ID")
	apiSecret := os.Getenv("APCA_API_SECRET_KEY")
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

func placeOrder(client alpaca.Client, db gorm.DB) (Order, error) {
	symbol := "BTC/USD"
	// qty := decimal.NewFromInt(1)
	qty := decimal.NewFromFloat(0.000038)
	side := alpaca.Side("buy")
	orderType := alpaca.OrderType("market")
	timeInForce := alpaca.TimeInForce("gtc") // day & ioc

	// Placing an order with the parameters set previously
	orderResult, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
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
		fmt.Printf("Order succesfully sent:\n%+v\n", *orderResult)
	}

	order := formatOrder(orderResult)

	saveOrder(order, db)

	return order, err
}

func formatOrder(orderResult *alpaca.Order) Order {

	order := Order{
		OrderId:        orderResult.ID,
		ClientOrderId:  orderResult.ClientOrderID,
		CreatedAt:      orderResult.CreatedAt,
		UpdatedAt:      orderResult.UpdatedAt,
		SubmittedAt:    orderResult.SubmittedAt,
		FilledAt:       orderResult.FilledAt,
		ExpiredAt:      orderResult.ExpiredAt,
		CanceledAt:     orderResult.CanceledAt,
		FailedAt:       orderResult.FilledAt,
		ReplacedAt:     orderResult.ReplacedAt,
		ReplacedBy:     "",
		Replaces:       "",
		AssetId:        orderResult.AssetID,
		Symbol:         orderResult.Symbol,
		AssetClass:     string(orderResult.OrderClass),
		Notional:       orderResult.Notional,
		Qty:            orderResult.Qty,
		FilledQty:      orderResult.FilledQty,
		FilledAvgPrice: orderResult.FilledAvgPrice,
		OrderType:      string(orderResult.Type),
		Type:           string(orderResult.Type),
		Side:           string(orderResult.Side),
		TimeInForce:    string(orderResult.TimeInForce),
		LimitPrice:     orderResult.LimitPrice,
		Status:         orderResult.Status,
	}

	return order
}

func saveOrder(order Order, db gorm.DB) (*gorm.DB, error) {
	// FirstOrCreate
	return db.Where(Order{OrderId: order.OrderId}).FirstOrCreate(&order), nil
}

func listPositions(client alpaca.Client, db gorm.DB) {
	// Get open positions
	positions, err := client.ListPositions()
	if err != nil {
		// Print error
		fmt.Printf("Failed to get open positions: %v\n", err)
	} else {
		// Print every position with its index
		for idx, position := range positions {
			fmt.Printf("Position %v: %+v\n", idx, position)
		}
	}
}
