package mono

import (
	"context"
	"net/http"
	"time"
)

// Unmarshaller allows to specify custom struct for unmarshalling response.
type Unmarshaller interface {
	Unmarshal(bts []byte, v interface{}) error
}

// HTTPClient defines which methods the library uses from `http.Client`.
// This also allows to unit-test the library.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Public is the client for accessing public API.
type Public interface {
	// Currency get basic list of currency.
	// The bank refreshes this list once in a five minutes or less.
	Currency(ctx context.Context) ([]CurrencyInfo, error)
}

// Personal is the client for accessing Personal API.
type Personal interface {
	// ClientInfo get info about the client for whom the token belongs.
	ClientInfo(ctx context.Context) (*UserInfo, error)
	Statements(ctx context.Context, account string, from, to time.Time) ([]StatementItem, error)
	LatestStatements(ctx context.Context, account string, from time.Time) ([]StatementItem, error)
}
