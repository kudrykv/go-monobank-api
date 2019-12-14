package mono

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type public struct {
	domain       string
	client       HTTPClient
	unmarshaller Unmarshaller
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

	if resp.StatusCode != http.StatusOK {
		var derp ErrorMono
		if err := p.unmarshaller.Unmarshal(bts, &derp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal body: %v", err)
		}

		return nil, fmt.Errorf("mono error: %s", derp.Description)
	}

	var currencies []CurrencyInfo
	if err := p.unmarshaller.Unmarshal(bts, &currencies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %v", err)
	}

	return currencies, nil
}
