package main

import (
	"context"
	"crypto-technologies/types"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata/stream"
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
	listenForMarket(db)
	listenForTrades()
}

func listenForMarket(db *gorm.DB) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Creating a client that connects to iex
	c := stream.NewCryptoClient(marketdata.US,
		stream.WithLogger(stream.DefaultLogger()),
		// configuring initial subscriptions and handlers
		stream.WithCryptoTrades(func(ct stream.CryptoTrade) {
			fmt.Printf("TRADE: %+v\n", ct)
		}, "*"),
		stream.WithCryptoQuotes(func(cq stream.CryptoQuote) {
			fmt.Printf("QUOTE: %+v\n", cq)

			rowsInserted, err := saveCryptoQuote(cq, *db)

			if rowsInserted != 1 || err != nil {
				panic(err)
			}

		}, "BTC/USD"),
		stream.WithCryptoOrderbooks(func(cob stream.CryptoOrderbook) {
			fmt.Printf("ORDERBOOK: %+v\n", cob)
		}, "BTC/USD"),
	)
	if err := c.Connect(ctx); err != nil {
		panic(err)
	}
	if err := <-c.Terminated(); err != nil {
		panic(err)
	}
}

func listenForTrades() {
	var tradeCount, quoteCount, barCount int32
	// modify these according to your needs
	tradeHandler := func(t stream.Trade) {
		atomic.AddInt32(&tradeCount, 1)
	}
	quoteHandler := func(q stream.Quote) {
		atomic.AddInt32(&quoteCount, 1)
	}
	barHandler := func(b stream.Bar) {
		atomic.AddInt32(&barCount, 1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setting up cancelling upon interrupt
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	go func() {
		<-s
		cancel()
	}()

	// Creating a client that connexts to iex
	c := stream.NewStocksClient(
		marketdata.IEX,
		// configuring initial subscriptions and handlers
		stream.WithTrades(tradeHandler, "BTC/USD"),
		stream.WithQuotes(quoteHandler, "BTC/USD", "ETH/USD"),
		stream.WithBars(barHandler, "BTC/USD", "ETH/USD"),
		// use stream.WithDailyBars to subscribe to daily bars too
		// use stream.WithCredentials to manually override envvars
		// use stream.WithHost to manually override envvar
		// use stream.WithLogger to use your own logger (i.e. zap, logrus) instead of log
		// use stream.WithProcessors to use multiple processing gourotines
		// use stream.WithBufferSize to change buffer size
		// use stream.WithReconnectSettings to change reconnect settings
	)

	// periodically displaying number of trades/quotes/bars received so far
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("trades:", tradeCount, "quotes:", quoteCount, "bars:", barCount)
		}
	}()

	if err := c.Connect(ctx); err != nil {
		log.Fatalf("could not establish connection, error: %s", err)
	}
	fmt.Println("established connection")

	// starting a goroutine that checks whether the client has terminated
	go func() {
		err := <-c.Terminated()
		if err != nil {
			log.Fatalf("terminated with error: %s", err)
		}
		fmt.Println("exiting")
		os.Exit(0)
	}()

	time.Sleep(3 * time.Second)
	// Adding BTC/USD trade subscription
	if err := c.SubscribeToTrades(tradeHandler, "BTC/USD"); err != nil {
		log.Fatalf("error during subscribing: %s", err)
	}
	fmt.Println("subscribed to BTC/USD trades")

	time.Sleep(3 * time.Second)
	// Unsubscribing from BTC/USD quotes
	if err := c.UnsubscribeFromQuotes("BTC/USD"); err != nil {
		log.Fatalf("error during unsubscribing: %s", err)
	}
	fmt.Println("unsubscribed from BTC/USD quotes")

	// and so on...
	time.Sleep(100 * time.Second)
	fmt.Println("we're done")
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
	db.AutoMigrate(&types.Account{}, &types.Order{}, &types.CryptoQuote{})

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
	qty := decimal.NewFromFloat(0.000038) // TODO: pass dynamic values
	tenCents := decimal.NewFromFloat(0.0000036)
	tpp := qty.Add(tenCents)
	spp := qty.Sub(tenCents.Mul(decimal.NewFromInt(2)))
	// limit := decimal.NewFromFloat(318.)
	tp := &alpaca.TakeProfit{LimitPrice: &tpp}
	sl := &alpaca.StopLoss{
		LimitPrice: nil,
		StopPrice:  &spp,
	}

	// qty := decimal.NewFromInt(1)

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
		OrderClass:  alpaca.Bracket,
		TakeProfit:  tp,
		StopLoss:    sl,
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

func saveCryptoQuote(marketQuote stream.CryptoQuote, db gorm.DB) (int64, error) {
	//types.CryptoQuote
	var cryptoQuote types.CryptoQuote = types.CryptoQuote{
		Symbol:    marketQuote.Symbol,
		Exchange:  marketQuote.Exchange,
		BidPrice:  marketQuote.BidPrice,
		BidSize:   marketQuote.BidSize,
		AskPrice:  marketQuote.AskPrice,
		AskSize:   marketQuote.AskSize,
		Timestamp: marketQuote.Timestamp,
	}

	result := db.Create(&cryptoQuote)
	return result.RowsAffected, result.Error
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
