package coinbase_test

import (
	"fmt"
	"os"

	"github.com/alpstable/coinbase"
)

func ExmapleClient_Accounts() {
	key := os.Getenv("COINBASE_API_KEY")
	secret := os.Getenv("COINBASE_API_SECRET")

	client, err := coinbase.NewClient(key, secret)
	if err != nil {
		panic(err)
	}

	accounts, err := client.Accounts()
	if err != nil {
		panic(err)
	}

	fmt.Println(accounts)
}
