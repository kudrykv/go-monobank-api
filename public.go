package mono

import (
	"context"
	"net/http"
)

type public struct {
	core
}

// NewPublic creates the client to interact with the public API.
//
// Factory supports parametrize http client, domain, and unmarshaller.
//  mono.NewPublic(
//    mono.WithDomain("https://custom-domain"),
//    mono.WithClient(&http.Client{}),
//    mono.WithUnmarshaller(&smthImplementingUnmarshaller),
//  )
func NewPublic(opts ...Option) Public {
	return public{core: newCore(opts...)}
}

func (p public) Currency(ctx context.Context) ([]CurrencyInfo, error) {
	var currencies []CurrencyInfo
	if err := p.request(ctx, http.MethodGet, p.domain+"/bank/currency", nil, &currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}
