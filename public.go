package mono

import (
	"context"
	"net/http"
)

type public struct {
	domain string

	tinyClient
}

func (p *public) setClient(client HTTPClient) {
	p.client = client
}

func (p *public) setDomain(domain string) {
	p.domain = domain
}

func (p *public) setUnmarshaller(u Unmarshaller) {
	p.unmarshaller = u
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
	p := public{}
	for _, opt := range opts {
		opt(&p)
	}

	if len(p.domain) == 0 {
		p.domain = "https://api.monobank.ua"
	}

	if p.client == nil {
		p.client = &http.Client{}
	}

	if p.unmarshaller == nil {
		p.unmarshaller = unmarshaller{}
	}

	return p
}

func (p public) Currency(ctx context.Context) ([]CurrencyInfo, error) {
	var currencies []CurrencyInfo
	if err := p.request(ctx, http.MethodGet, p.domain+"/bank/currency", nil, &currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}
