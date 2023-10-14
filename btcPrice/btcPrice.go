package btcPrice

import (
	"context"
	"crypto-technologies/types"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	ChannelTicker   string = "ticker"
	TypeSubscribe   string = "subscribe"
	TypeUnsubscribe string = "unsubscribe"
)

const Address = "wss://ws-feed.exchange.coinbase.com"

type Message struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type Trade struct {
	Type    string    `json:"type"`
	Time    time.Time `json:"time"`
	Product string    `json:"product_id"`
	Price   string    `json:"price"`
}

func GetBtcPrice(db gorm.DB) {
	// 1. websocket client connection
	conn, _, err := websocket.Dial(context.Background(), Address, nil)

	if err != nil {
		panic(err)
	}

	defer conn.Close(http.StatusOK, "connection closed")

	// 2. subscribe
	sub := Message{
		Type:       TypeSubscribe,
		ProductIDs: []string{"BTC-USD"},
		Channels:   []string{ChannelTicker},
	}

	err = wsjson.Write(context.Background(), conn, sub)

	if err != nil {
		log.Error().Err(fmt.Errorf("failed to subscribe to coinbase channel '%s': %w", "BTC-USD", err))
		return
	}

	// 3. read from the websocket
	for {
		_, message, err := conn.Read(context.Background())

		if err != nil {
			break
		}

		var trade Trade

		json.Unmarshal(message, &trade)

		fmt.Println(trade.Price)
		saveBtcPrice(trade, db)

	}
}

func saveBtcPrice(btcPrice Trade, db gorm.DB) (int64, error) {
	var currentPrice types.BtcPrice = types.BtcPrice{
		Type:    btcPrice.Type,
		Time:    btcPrice.Time,
		Product: btcPrice.Product,
		Price:   btcPrice.Price,
	}

	result := db.Create(&currentPrice)
	return result.RowsAffected, result.Error
}
