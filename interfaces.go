package mono

import (
	"context"
	"net/http"
)

type Unmarshaller interface {
	Unmarshal(bts []byte, v interface{}) error
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Public interface {
	Currency(ctx context.Context) ([]CurrencyInfo, error)
}
