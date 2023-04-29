package coinbase_test

import (
	"context"
	"log"
	"os"

	"github.com/alpstable/coinbase"
	"github.com/google/uuid"
)

//nolint:testableexamples
func ExampleClient_Accounts() {
	// DO NOT RUN THIS EXAMPLE ON A LIVE ACCOUNT
	key := os.Getenv("COINBASE_API_KEY")
	secret := os.Getenv("COINBASE_API_SECRET")

	if key == "" || secret == "" {
		log.Println("COINBASE_API_KEY or COINBASE_API_SECRET is empty")

		return
	}

	client, err := coinbase.NewClient(key, secret)
	if err != nil {
		panic(err)
	}

	accounts, err := client.Accounts(context.Background())
	if err != nil {
		panic(err)
	}

	log.Printf("accounts: %+v", accounts)
}

//nolint:testableexamples
func ExampleClient_CreateOrder() {
	// DO NOT RUN THIS EXAMPLE ON A LIVE ACCOUNT
	key := os.Getenv("COINBASE_API_KEY")
	secret := os.Getenv("COINBASE_API_SECRET")

	if key == "" || secret == "" {
		log.Println("COINBASE_API_KEY or COINBASE_API_SECRET is empty")

		return
	}

	client, err := coinbase.NewClient(key, secret)
	if err != nil {
		panic(err)
	}

	limit := &coinbase.LimitGTCConfig{
		BaseSize: "1",
		Price:    "2.7",
	}

	confg := coinbase.OrderConfig{
		LimitGTC: limit,
	}

	req := coinbase.OrderRequest{
		ClientOrderID: uuid.New().String(),
		ProductID:     "BTC-USDT",
		Side:          coinbase.OrderSideBuy,
		Configuration: confg,
	}

	// Buy 1 BTC at 2.7 USDT
	orderResponse, err := client.CreateOrder(context.Background(), req)
	if err != nil {
		panic(err)
	}

	log.Printf("orderResponse: %+v", orderResponse)
}
