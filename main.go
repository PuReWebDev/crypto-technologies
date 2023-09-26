package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/joho/godotenv"
	// "github.com/shopspring/decimal"	// Used for decimal type
)

func getAccount(client alpaca.Client) {
	// Get account information
	acct, err := client.GetAccount()
	if err != nil {
		// Print error
		fmt.Printf("Error getting account information: %v", err)
	} else {
		// Print account information
		fmt.Printf("Account: %+v\n", *acct)
	}
}

func main() {

	loadEnvironment()

	client := prepAlpaca()

	getAccount(client)
}

func loadEnvironment() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
