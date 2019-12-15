package mono

import (
	"context"
	"net/http"
)

type personal struct {
	core
}

// NewPersonal creates the client to access Personal API.
func NewPersonal(apiToken string, opts ...Option) Personal {
	if len(apiToken) == 0 {
		panic("api token is required")
	}

	p := personal{core: newCore(opts...)}
	p.token = apiToken

	return p
}

func (p personal) ClientInfo(ctx context.Context) (*UserInfo, error) {
	var userInfo UserInfo
	if err := p.request(ctx, http.MethodGet, p.domain+"/personal/client-info", nil, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
