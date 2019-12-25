package mono

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func (p personal) LatestStatements(ctx context.Context, account string, from time.Time) ([]StatementItem, error) {
	return p.Statements(ctx, account, from, time.Now())
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

func (p personal) SetWebhook(ctx context.Context, webhook string) error {
	wh := strings.ReplaceAll(strings.ReplaceAll(webhook, `"`, `\"`), "\n", "\\n")
	body := bytes.NewReader([]byte(`{"webHookUrl":"` + wh + `"}`))

	var empty struct{}
	return p.request(ctx, http.MethodPost, p.domain+"/personal/webhook", body, &empty)
}

func (p personal) ParseWebhook(_ context.Context, rc io.ReadCloser) (*WebhookData, error) {
	bts, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	if err := rc.Close(); err != nil {
		return nil, fmt.Errorf("failed to close body: %v", err)
	}

	var wh WebhookData
	if err := p.unmarshaller.Unmarshal(bts, &wh); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bytes: %v", err)
	}

	return &wh, nil
}

func (p personal) ListenForWebhooks(_ context.Context) (<-chan WebhookData, http.HandlerFunc) {
	whch := make(chan WebhookData, p.whBufferSize)

	return whch, func(w http.ResponseWriter, r *http.Request) {
		wh, err := p.ParseWebhook(r.Context(), r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go func() {
			whch <- *wh
		}()

		w.WriteHeader(http.StatusOK)
	}
}
