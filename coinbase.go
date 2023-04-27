// coinbase is a Go client library for the Coinbase API.

package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const api = "https://api.coinbase.com/api/v3"

// Client is a Coinbase API client.
type Client struct {
	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}
}

// NewClient creates a new Coinbase API client with the provided API key and
// secret. The Coinbase API requests are automatically signed with the provided
// API key and secret using an http Transport middleware.
func NewClient(key, secret string) (*Client, error) {
	httpClient := http.DefaultClient

	var err error

	httpClient.Transport, err = newRoundTripper(key, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	client := &Client{
		httpClient: http.DefaultClient,
	}

	return client, nil
}

// AvailableMoney is the amount of money that is available to send.
type AvailableMoney struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// HoldMoney is the amount of money that is on hold.
type HoldMoney struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// Account is an authenticated user's account.
type Account struct {
	UUID             string         `json:"uuid"`
	Name             string         `json:"name"`
	Currency         string         `json:"currency"`
	AvailableBalance AvailableMoney `json:"available_balance"`
	Default          bool           `json:"default"`
	Active           bool           `json:"active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        *time.Time     `json:"deleted_at,omitempty"`
	Type             string         `json:"type"`
	Ready            bool           `json:"ready"`
	Hold             HoldMoney      `json:"hold"`
}

// Accounts is a list of authenticated accounts for the current user.
type Accounts struct {
	Data    []Account `json:"accounts"`
	HasNext bool      `json:"has_next"`
	Cursor  string    `json:"cursor"`
	Size    int32     `json:"size"`
}

// Accounts returns a slice of accounts for the authenticated user.
//
// https://docs.cloud.coinbase.com/advanced-trade-api/reference/retailbrokerageapi_getaccounts
func (client *Client) Accounts(ctx context.Context) (*Accounts, error) {
	full, err := url.JoinPath(api, "brokerage", "accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	accounts := &Accounts{}
	if err := json.NewDecoder(resp.Body).Decode(accounts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return accounts, nil
}
