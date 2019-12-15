package mono

import (
	"context"
	"net/http"
)

type personal struct {
	token        string
	domain       string
	client       HTTPClient
	unmarshaller Unmarshaller
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

func NewPersonal(apiToken string, opts ...Option) Personal {
	if len(apiToken) == 0 {
		panic("api token is required")
	}

	p := personal{token: apiToken}

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

	err := tinyClient{client: p.client, unmarshaller: p.unmarshaller, token: p.token}.
		request(ctx, http.MethodGet, p.domain+"/personal/client-info", nil, &userInfo)

	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
