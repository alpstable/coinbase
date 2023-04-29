package coinbase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type mockClient struct {
	response   []byte
	statusCode int
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer(m.response)),
		StatusCode: m.statusCode,
	}, nil
}

func TestAccounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response []byte
		want     *Accounts
		err      error
	}{
		{
			name: "nil",
			err:  io.EOF, // end of file, nothing in response
		},
		{
			name:     "empty slice",
			response: []byte(`{}`),
			want:     &Accounts{},
		},
		{
			name: "single",
			response: []byte(`
{
  "accounts": [{
    "uuid": "8bfc20d7-f7c6-4422-bf07-8243ca4169fe",
    "name": "BTC Wallet",
    "currency": "BTC",
    "available_balance": {
      "value": "1.23",
      "currency": "BTC"
    },
    "default": false,
    "active": true,
    "created_at": "2021-05-31T09:59:59Z",
    "updated_at": "2021-05-31T09:59:59Z",
    "deleted_at": "2021-05-31T09:59:59Z",
    "type": "ACCOUNT_TYPE_UNSPECIFIED",
    "ready": true,
    "hold": {
      "value": "1.23",
      "currency": "BTC"
    }
  }],
  "has_next": true,
  "cursor": "789100",
  "size": 1
}`),
			want: &Accounts{
				Data: []Account{
					{
						UUID:     "8bfc20d7-f7c6-4422-bf07-8243ca4169fe",
						Name:     "BTC Wallet",
						Currency: "BTC",
						AvailableBalance: AvailableMoney{
							Value:    "1.23",
							Currency: "BTC",
						},
						Default:   false,
						Active:    true,
						CreatedAt: time.Date(2021, 5, 31, 9, 59, 59, 0, time.UTC),
						UpdatedAt: time.Date(2021, 5, 31, 9, 59, 59, 0, time.UTC),
						DeletedAt: func() *time.Time {
							dt := time.Date(2021, 5, 31, 9, 59, 59, 0, time.UTC)

							return &dt
						}(),
						Type:  "ACCOUNT_TYPE_UNSPECIFIED",
						Ready: true,
						Hold: HoldMoney{
							Value:    "1.23",
							Currency: "BTC",
						},
					},
				},
				HasNext: true,
				Cursor:  "789100",
				Size:    1,
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &Client{
				httpClient: &mockClient{
					response:   test.response,
					statusCode: http.StatusOK,
				},
			}

			got, err := client.Accounts(context.Background())
			if !errors.Is(err, test.err) {
				t.Fatalf("got %v, want %v", err, test.err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response []byte
		want     *Order
		err      error
	}{
		{
			name: "nil",
			err:  io.EOF, // end of file, nothing in response
		},
		{
			name:     "empty slice",
			response: []byte(`{}`),
			want:     &Order{},
		},
		{
			name: "single",
			response: []byte(`
{
  "success": true,
  "failure_reason": "string",
  "order_id": "string",
  "success_response": {
    "order_id": "11111-00000-000000",
    "product_id": "BTC-USD",
    "side": "UNKNOWN_ORDER_SIDE",
    "client_order_id": "0000-00000-000000"
  },
  "error_response": {
    "error": "UNKNOWN_FAILURE_REASON",
    "message": "The order configuration was invalid",
    "error_details": "Market orders cannot be placed with empty order sizes",
    "preview_failure_reason": "UNKNOWN_PREVIEW_FAILURE_REASON",
    "new_order_failure_reason": "UNKNOWN_FAILURE_REASON"
  },
  "order_configuration": {
    "market_market_ioc": {
      "quote_size": "10.00",
      "base_size": "0.001"
    },
    "limit_limit_gtc": {
      "base_size": "0.001",
      "limit_price": "10000.00",
      "post_only": false
    },
    "limit_limit_gtd": {
      "base_size": "0.001",
      "limit_price": "10000.00",
      "end_time": "2021-05-31T09:59:59Z",
      "post_only": false
    },
    "stop_limit_stop_limit_gtc": {
      "base_size": "0.001",
      "limit_price": "10000.00",
      "stop_price": "20000.00",
      "stop_direction": "UNKNOWN_STOP_DIRECTION"
    },
    "stop_limit_stop_limit_gtd": {
      "base_size": "0.001",
      "limit_price": "10000.00",
      "stop_price": "20000.00",
      "end_time": "2021-05-31T09:59:59Z",
      "stop_direction": "UNKNOWN_STOP_DIRECTION"
    }
  }
}
`),
			want: &Order{
				Success:       true,
				FailureReason: "string",
				OrderID:       "string",
				SuccessResponse: SuccessResponse{
					OrderID:       "11111-00000-000000",
					ProductID:     "BTC-USD",
					Side:          "UNKNOWN_ORDER_SIDE",
					ClientOrderID: "0000-00000-000000",
				},
				ErrorResponse: ErrorResponse{
					Error:                 "UNKNOWN_FAILURE_REASON",
					Message:               "The order configuration was invalid",
					ErrorDetails:          "Market orders cannot be placed with empty order sizes",
					PreviewFailureReason:  "UNKNOWN_PREVIEW_FAILURE_REASON",
					NewOrderFailureReason: "UNKNOWN_FAILURE_REASON",
				},
				OrderConfiguration: OrderConfig{
					MarketIOC: &MarketIOCConfig{
						QuoteSize: "10.00",
						BaseSize:  "0.001",
					},
					LimitGTC: &LimitGTCConfig{
						BaseSize: "0.001",
						Price:    "10000.00",
						PostOnly: false,
					},
					LimitGTD: &LimitGTDConfig{
						BaseSize: "0.001",
						Price:    "10000.00",
						EndTime:  time.Date(2021, 5, 31, 9, 59, 59, 0, time.UTC),
						PostOnly: false,
					},
					StopLimitGTC: &StopLimitGTCConfig{
						BaseSize:      "0.001",
						LimitPrice:    "10000.00",
						StopPrice:     "20000.00",
						StopDirection: "UNKNOWN_STOP_DIRECTION",
					},
					StopLimitGTD: &StopLimitGTDConfig{
						BaseSize:      "0.001",
						LimitPrice:    "10000.00",
						StopPrice:     "20000.00",
						EndTime:       time.Date(2021, 5, 31, 9, 59, 59, 0, time.UTC),
						StopDirection: "UNKNOWN_STOP_DIRECTION",
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &Client{
				httpClient: &mockClient{
					response:   test.response,
					statusCode: http.StatusOK,
				},
			}

			got, err := client.CreateOrder(context.Background(), OrderRequest{})
			if !errors.Is(err, test.err) {
				t.Fatalf("got %v, want %v", err, test.err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}
