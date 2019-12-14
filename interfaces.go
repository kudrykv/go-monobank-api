package mono

import (
	"context"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Public interface {
	Currency(ctx context.Context) ([]CurrencyInfo, error)
}
