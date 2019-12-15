package mono

import (
	"context"
	"net/http"
)

type personal struct {
	domain string

	tinyClient
}

func (p *personal) setClient(client HTTPClient) {
	p.client = client
}

func (p *personal) setDomain(domain string) {
	p.domain = domain
}

func (p *personal) setUnmarshaller(u Unmarshaller) {
	p.unmarshaller = u
}

// NewPersonal creates the client to access Personal API.
func NewPersonal(apiToken string, opts ...Option) Personal {
	if len(apiToken) == 0 {
		panic("api token is required")
	}

	p := personal{tinyClient: tinyClient{token: apiToken}}

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

func (p personal) ClientInfo(ctx context.Context) (*UserInfo, error) {
	var userInfo UserInfo
	if err := p.request(ctx, http.MethodGet, p.domain+"/personal/client-info", nil, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
