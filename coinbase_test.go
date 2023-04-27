package coinbase

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type mockClient struct {
	response []byte
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: io.NopCloser(bytes.NewBuffer(m.response)),
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
					response: test.response,
				},
			}

			got, err := client.Accounts()
			if !errors.Is(err, test.err) {
				t.Fatalf("got %v, want %v", err, test.err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}
