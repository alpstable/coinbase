package coinbase_test

import (
	"context"
	"log"
	"os"

	"github.com/alpstable/coinbase"
)

func ExampleClient_Accounts() {
	key := os.Getenv("COINBASE_API_KEY")
	secret := os.Getenv("COINBASE_API_SECRET")

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
