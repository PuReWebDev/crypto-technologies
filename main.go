package main

import (
	"crypto-technologies/types"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getAccount(client *alpaca.Client, db gorm.DB) {
	// Get account information
	acct, err := client.GetAccount()

	account, error := buildAccount(acct)

	if error != nil {
		fmt.Printf("Error: not able to build the account: %v\n", error)
	}

	if err != nil {
		// Print error
		fmt.Printf("Error getting account information: %v", err)
	} else {
		// Print account information
		fmt.Printf("Account ID: %+v\n", acct.ID)
		fmt.Printf("Account: %+v\n", *acct)

		// FirstOrCreate
		result := db.Where(types.Account{AccountID: account.AccountID}).FirstOrCreate(&account)

		fmt.Printf("Query Result: %+v\n", result)
	}
}

func main() {

	db := loadEnvironment()

	client := prepAlpaca()

	getAccount(client, *db)
	placeOrder(client, *db)
	listPositions(client, *db)
	getAsset(client)
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
	db.AutoMigrate(&types.Account{}, &types.Order{})

	return db
}

func prepAlpaca() *alpaca.Client {
	apiKey := os.Getenv("APCA_API_KEY_ID")
	apiSecret := os.Getenv("APCA_API_SECRET_KEY")
	baseURL := os.Getenv("baseURL")

	client := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseURL,
	})

	return client
}

func placeOrder(client *alpaca.Client, db gorm.DB) (types.Order, error) {
	symbol := "BTC/USD"
	// qty := decimal.NewFromInt(1)
	qty := decimal.NewFromFloat(0.000038) // TODO: pass dynamic values
	side := alpaca.Side("buy")
	orderType := alpaca.OrderType("market")
	timeInForce := alpaca.TimeInForce("gtc") // day & ioc

	// Placing an order with the parameters set previously
	orderResult, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
		Symbol:      symbol,
		Qty:         &qty,
		Side:        side,
		Type:        orderType,
		TimeInForce: timeInForce,
	})

	var order types.Order

	if err != nil {
		// Print error
		fmt.Printf("Failed to place order: %v\n", err)
	} else {
		// Print resulting order object
		fmt.Printf("Order succesfully sent:\n%+v\n", *orderResult)

		order := formatOrder(orderResult)

		saveOrder(order, db)
	}

	return order, err
}

func formatOrder(orderResult *alpaca.Order) types.Order {

	order := types.Order{
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
		ReplacedBy:     *orderResult.ReplacedBy,
		Replaces:       *orderResult.Replaces,
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

func buildAccount(acct *alpaca.Account) (types.Account, error) {

	var account types.Account = types.Account{
		AccountID:             acct.ID,
		AccountBlocked:        acct.AccountBlocked,
		AccountNumber:         acct.AccountNumber,
		PortfolioValue:        acct.PortfolioValue,
		Status:                acct.Status,
		Currency:              acct.Currency,
		Cash:                  acct.Cash,
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
		CryptoStatus:          acct.CryptoStatus,
		DaytradeCount:         acct.DaytradeCount,
		LastMaintenanceMargin: acct.LastMaintenanceMargin,
		DaytradingBuyingPower: acct.DaytradingBuyingPower,
		RegtBuyingPower:       acct.RegTBuyingPower,
		CreatedAt:             acct.CreatedAt,
		UpdatedAt:             time.Now(),
	}

	return account, nil
}

func saveOrder(order types.Order, db gorm.DB) (*gorm.DB, error) {
	// FirstOrCreate
	return db.Where(types.Order{OrderId: order.OrderId}).FirstOrCreate(&order), nil
}

func listPositions(client *alpaca.Client, db gorm.DB) {
	// Get the last 100 of our closed orders
	status := "closed"
	limit := 100
	nested := true // show nested multi-leg orders

	var orderRequest alpaca.GetOrdersRequest = alpaca.GetOrdersRequest{
		Status: status,
		Limit:  limit,
		Nested: nested,
	}

	positions, err := alpaca.GetOrders(orderRequest)

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

func getAsset(client *alpaca.Client) (*alpaca.Asset, error) {
	res, error := client.GetAsset("BTC")

	if error != nil {
		fmt.Printf("Failed to get asset data: %v\n", error)
	} else {
		fmt.Printf("Success, Asset Data: %v\n", res)
	}

	return res, nil
}
