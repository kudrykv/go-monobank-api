package mono

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Public interface {
	Currency(ctx context.Context) ([]CurrencyInfo, error)
}

type public struct {
	domain string
	client HTTPClient
}

func (p *public) setClient(client HTTPClient) {
	p.client = client
}

func (p *public) setDomain(domain string) {
	p.domain = domain
}

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

	return p
}

func (p public) Currency(ctx context.Context) ([]CurrencyInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.domain+"/bank/currency", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the body: %v", err)
	}

	var currencies []CurrencyInfo
	if err := json.Unmarshal(bts, &currencies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %v", err)
	}

	return currencies, nil
}
