package mono

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"
)

const MaxAllowedDuration = 2682000

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

func (p personal) Statements(ctx context.Context, account string, from, to time.Time) ([]StatementItem, error) {
	if len(account) == 0 {
		return nil, errors.New("account must be set")
	}

	if from.After(to) {
		return nil, errors.New("`from` should be less than `to`")
	}

	if to.Unix()-from.Unix() > MaxAllowedDuration {
		return nil, errors.New("max allowed duration is " + strconv.Itoa(MaxAllowedDuration) + " seconds")
	}

	fromUnix := strconv.FormatInt(from.Unix(), 10)
	toUnix := strconv.FormatInt(to.Unix(), 10)
	url := p.domain + "/personal/statement/" + account + "/" + fromUnix + "/" + toUnix

	var statements []StatementItem
	if err := p.request(ctx, http.MethodGet, url, nil, &statements); err != nil {
		return nil, err
	}

	return statements, nil
}
