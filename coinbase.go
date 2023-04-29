// coinbase is a Go client library for the Coinbase API.

package coinbase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const api = "https://api.coinbase.com/api/v3"

// ErrStatusNotOK is returned when the Coinbase API returns a non-OK status
// code.
var ErrStatusNotOK = errors.New("status not OK")

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

// AvailableMoney represents an amount of money that is available.
type AvailableMoney struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// HoldMoney represents an amount of money that is being held.
type HoldMoney struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// Account represents a user account with the available balance and hold amount
// of currency.
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

// Accounts represents a collection of accounts along with metadata.
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("%w: unexpected status code: %d, body: %s",
			ErrStatusNotOK, resp.StatusCode, body)
	}

	accounts := &Accounts{}
	if err := json.NewDecoder(resp.Body).Decode(accounts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return accounts, nil
}

// MarketIOCConfig represents the configuration of a market or
// immediate-or-cancel order.
type MarketIOCConfig struct {
	QuoteSize string `json:"quote_size" validate:"required_if=Side:BUY"`
	BaseSize  string `json:"base_size" validate:"required_if=Side:SELL"`
}

// LimitGTCConfig represents the configuration of a good-'til-cancelled limit
// order.
type LimitGTCConfig struct {
	BaseSize string `json:"base_size" validate:"required"`
	Price    string `json:"limit_price" validate:"required"`
	PostOnly bool   `json:"post_only"`
}

// LimitGTDConfig represents the configuration of a good-'til-date limit order.
type LimitGTDConfig struct {
	BaseSize string    `json:"base_size" validate:"required"`
	Price    string    `json:"limit_price" validate:"required"`
	EndTime  time.Time `json:"end_time" validate:"required"`
	PostOnly bool      `json:"post_only"`
}

// OrderStopDirection represents the possible stop directions for an order.
type OrderStopDirection string

const (
	// StopDirUnknown represents an unknown stop direction for an order.
	StopDirUnknown OrderStopDirection = "UNKNOWN_STOP_DIRECTION"

	// StopDirUp represents a stop direction for an order that triggers if
	// the price goes up.
	StopDirUp OrderStopDirection = "STOP_DIRECTION_STOP_UP"

	// StopDirDown represents a stop direction for an order that triggers if
	// the price goes down.
	StopDirDown OrderStopDirection = "STOP_DIRECTION_STOP_DOWN"
)

// StopLimitGTCConfig represents a stop-limit order with Good-til-Canceled time
// in force.
type StopLimitGTCConfig struct {
	BaseSize      string             `json:"base_size" validate:"required"`
	LimitPrice    string             `json:"limit_price" validate:"required"`
	StopPrice     string             `json:"stop_price" validate:"required"`
	StopDirection OrderStopDirection `json:"stop_direction"`
}

// StopLimitGTDConfig represents a stop-limit order with Good-til-Date time in
// force.
type StopLimitGTDConfig struct {
	BaseSize      string             `json:"base_size" validate:"required"`
	LimitPrice    string             `json:"limit_price" validate:"required"`
	StopPrice     string             `json:"stop_price" validate:"required"`
	StopDirection OrderStopDirection `json:"stop_direction"`
	EndTime       time.Time          `json:"end_time" validate:"required"`
}

// OrderConfig represents the configuration of an order.
type OrderConfig struct {
	MarketIOC    *MarketIOCConfig    `json:"market_market_ioc,omitempty"`
	LimitGTC     *LimitGTCConfig     `json:"limit_limit_gtc,omitempty"`
	LimitGTD     *LimitGTDConfig     `json:"limit_limit_gtd,omitempty"`
	StopLimitGTC *StopLimitGTCConfig `json:"stop_limit_stop_limit_gtc,omitempty"`
	StopLimitGTD *StopLimitGTDConfig `json:"stop_limit_stop_limit_gtd,omitempty"`
}

// OrderSide represents the side of an order, either BUY or SELL.
type OrderSide string

const (
	// OrderSideUnknown represents an unknown or undefined order side.
	OrderSideUnknown OrderSide = "UNKNOWN_ORDER_SIDE"

	// OrderSideBuy represents a buy order side.
	OrderSideBuy OrderSide = "BUY"

	// OrderSideSell represents a sell order side.
	OrderSideSell OrderSide = "SELL"
)

// OrderRequest can be used to create an order on Coinbase.
type OrderRequest struct {
	ClientOrderID string      `json:"client_order_id" validate:"required"`
	ProductID     string      `json:"product_id" validate:"required"`
	Side          OrderSide   `json:"side"`
	Configuration OrderConfig `json:"order_configuration"`
}

// SuccessResponse represents a successful order response.
type SuccessResponse struct {
	OrderID       string    `json:"order_id"`
	ProductID     string    `json:"product_id,omitempty"`
	Side          OrderSide `json:"side,omitempty"`
	ClientOrderID string    `json:"client_order_id,omitempty"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error                 string `json:"error"`
	Message               string `json:"message,omitempty"`
	ErrorDetails          string `json:"error_details,omitempty"`
	PreviewFailureReason  string `json:"preview_failure_reason,omitempty"`
	NewOrderFailureReason string `json:"new_order_failure_reason,omitempty"`
}

// Order is the response from creating an order.
type Order struct {
	Success            bool            `json:"success"`
	FailureReason      string          `json:"failure_reason"`
	OrderID            string          `json:"order_id"`
	SuccessResponse    SuccessResponse `json:"success_response,omitempty"`
	ErrorResponse      ErrorResponse   `json:"error_response,omitempty"`
	OrderConfiguration OrderConfig     `json:"order_configuration,omitempty"`
}

// CreateOrder will create an order with a specified product_id (BASE-QUOTE),
// side (buy/sell), etc.
//
// https://docs.cloud.coinbase.com/advanced-trade-api/reference/retailbrokerageapi_postorder
func (client *Client) CreateOrder(ctx context.Context, orderReq OrderRequest) (*Order, error) {
	full, err := url.JoinPath(api, "brokerage", "orders")
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	// Create the request body.
	body, err := json.Marshal(orderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, full, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Header should be application/json.
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("%w: unexpected status code: %d, body: %s",
			ErrStatusNotOK, resp.StatusCode, body)
	}

	orderResponse := &Order{}
	if err := json.NewDecoder(resp.Body).Decode(orderResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return orderResponse, nil
}
