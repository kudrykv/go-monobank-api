package mono_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mono "github.com/kudrykv/go-monobank-api"
)

var personalResponseBody = `{
  "name": "deadbeef",
  "webHookUrl": "https://url/leading/to/the/webhook",
  "accounts": [
    {
      "id": "kKGVoZuHWzqVoZuH",
      "balance": 10000000,
      "creditLimit": 10000000,
      "currencyCode": 980,
      "cashbackType": "UAH"
    }
  ]
}`

var expectedPersonal = mono.UserInfo{
	Name:       "deadbeef",
	WebHookURL: "https://url/leading/to/the/webhook",
	Accounts: []mono.Account{{
		ID:                  "kKGVoZuHWzqVoZuH",
		Balance:             10000000,
		CreditLimit:         10000000,
		CurrencyCodeISO4217: 980,
		CashbackType:        mono.CashbackUAH,
	}},
}

func TestNewPersonal(t *testing.T) {
	p := mono.NewPersonal("token")
	expectNotNil(t, p)
}

func TestNewPersonal_PanicOnEmptyToken(t *testing.T) {
	defer func() {
		err, ok := recover().(string)

		expectTrue(t, ok)
		expectNotNil(t, err)
		expectError(t, errors.New(err), "api token is required")
	}()

	mono.NewPersonal("")
}

func TestPersonal_ClientInfo_Succ(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(personalResponseBody))),
	}

	ctx := context.Background()
	apiToken := "api-token"

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	actual, err := personal.ClientInfo(ctx)

	expectNoError(t, err)
	expectDeepEquals(t, actual, &expectedPersonal)
}

func TestPersonal_ClientInfo_Fail(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(failResponseBody))),
	}

	ctx := context.Background()
	apiToken := "api-token"

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	_, err := personal.ClientInfo(ctx)
	expectError(t, err, "mono error: go away")
}

var statementsResponseBody = `[
  {
    "id": "ZuHWzqkKGVo=",
    "time": 1554466347,
    "description": "Покупка щастя",
    "mcc": 7997,
    "hold": false,
    "amount": -95000,
    "operationAmount": -95000,
    "currencyCode": 980,
    "commissionRate": 0,
    "cashbackAmount": 19000,
    "balance": 10050000
  }
]`

var expectedStatementsResponse = []mono.StatementItem{{
	ID:                  "ZuHWzqkKGVo=",
	Time:                1554466347,
	Description:         "Покупка щастя",
	MCC:                 7997,
	Hold:                false,
	Amount:              -95000,
	OperationAmount:     -95000,
	CurrencyCodeISO4217: 980,
	CommissionRate:      0,
	CashbackAmount:      19000,
	Balance:             10050000,
}}

func TestPersonal_Statements_Succ(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(statementsResponseBody))),
	}

	ctx := context.Background()
	apiToken := "api-token"
	accountID := "deadbeef"
	from := time.Now().Add(-time.Hour * 24 * 15) // 15 days
	to := time.Now()

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	statements, err := personal.Statements(ctx, accountID, from, to)
	expectNoError(t, err)
	expectDeepEquals(t, statements, expectedStatementsResponse)
}

func TestPersonal_Statements_Fail(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(failResponseBody))),
	}

	ctx := context.Background()
	apiToken := "api-token"
	accountID := "deadbeef"
	from := time.Now().Add(-time.Hour * 24 * 15) // 15 days
	to := time.Now()

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	_, err := personal.Statements(ctx, "", from, to)
	expectError(t, err, "account must be set")

	_, err = personal.Statements(ctx, accountID, to, from)
	expectError(t, err, "`from` should be less than `to`")

	_, err = personal.Statements(ctx, accountID, from.Add(-time.Hour*24*45), to)
	expectError(t, err, "max allowed duration is 2682000 seconds")

	_, err = personal.Statements(ctx, accountID, from, to)
	expectError(t, err, "mono error: go away")
}

func TestPersonal_LatestStatements(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(statementsResponseBody))),
	}

	ctx := context.Background()
	apiToken := "api-token"
	accountID := "deadbeef"
	from := time.Now().Add(-time.Hour * 24 * 15) // 15 days

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	statements, err := personal.LatestStatements(ctx, accountID, from)
	expectNoError(t, err)
	expectDeepEquals(t, statements, expectedStatementsResponse)
}

func TestPersonal_SetWebhook_Succ(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"status":"ok"}`))),
	}

	ctx := context.Background()
	apiToken := "api-token"
	wh := "https://domain/webhook"

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	err := personal.SetWebhook(ctx, wh)
	expectNoError(t, err)
}

func TestPersonal_SetWebhook_Fail(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(failResponseBody))),
	}

	ctx := context.Background()
	apiToken := "api-token"
	wh := "https://domain/webhook"

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	err := personal.SetWebhook(ctx, wh)
	expectError(t, err, "mono error: go away")
}

var webhookBody = `{
	"type": "StatementItem",
	"data": {
		"account": "deadbeef",
		"statementItem": {
			"id": "ZuHWzqkKGVo=",
			"time": 1554466347,
			"description": "Покупка щастя",
			"mcc": 7997,
			"hold": false,
			"amount": -95000,
			"operationAmount": -95000,
			"currencyCode": 980,
			"commissionRate": 0,
			"cashbackAmount": 19000,
			"balance": 10050000
		}
	}
}`

var webhookParsed = mono.WebhookData{
	Type: "StatementItem",
	Data: mono.WebhookStatementItem{
		AccountID: "deadbeef",
		StatementItem: mono.StatementItem{
			ID:                  "ZuHWzqkKGVo=",
			Time:                1554466347,
			Description:         "Покупка щастя",
			MCC:                 7997,
			Hold:                false,
			Amount:              -95000,
			OperationAmount:     -95000,
			CurrencyCodeISO4217: 980,
			CommissionRate:      0,
			CashbackAmount:      19000,
			Balance:             10050000,
		},
	},
}

func TestPersonal_ParseWebhook_Succ(t *testing.T) {
	client := &clienttest{}

	ctx := context.Background()
	apiToken := "api-token"
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(webhookBody)))

	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	actual, err := personal.ParseWebhook(ctx, req.Body)
	expectNoError(t, err)
	expectDeepEquals(t, actual, &webhookParsed)
}

func TestPersonal_ParseWebhook_Fail(t *testing.T) {
	um := &unmtest{}
	client := &clienttest{}

	ctx := context.Background()
	apiToken := "api-token"
	personal := mono.NewPersonal(apiToken, mono.WithClient(client), mono.WithUnmarshaller(um))

	req := httptest.NewRequest(http.MethodGet, "/", &badReader{})

	_, err := personal.ParseWebhook(ctx, req.Body)
	expectErrorStartsWith(t, err, "failed to read body: ")

	req = httptest.NewRequest(http.MethodGet, "/", &badReadCloser{})

	_, err = personal.ParseWebhook(ctx, req.Body)
	expectErrorStartsWith(t, err, "failed to close body: ")

	req = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(webhookBody)))
	um.Err = errors.New("boo")

	_, err = personal.ParseWebhook(ctx, req.Body)
	expectErrorStartsWith(t, err, "failed to unmarshal bytes: ")
}

func TestPersonal_ListenForWebhooks_Succ(t *testing.T) {
	client := &clienttest{}

	apiToken := "api-token"
	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	whChan, handlerFunc := personal.ListenForWebhooks(context.Background())

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(webhookBody)))

	handlerFunc(w, r)
	expectEquals(t, w.Code, http.StatusOK)

	select {
	case <-time.After(time.Millisecond * 10):
		t.Fatalf("died waiting on the message")

	case wh := <-whChan:
		expectDeepEquals(t, wh, webhookParsed)
	}
}

func TestPersonal_ListenForWebhooks_Fail(t *testing.T) {
	client := &clienttest{}

	apiToken := "api-token"
	personal := mono.NewPersonal(apiToken, mono.WithClient(client))

	_, handlerFunc := personal.ListenForWebhooks(context.Background())

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", &badReadCloser{})

	handlerFunc(w, r)
	expectEquals(t, w.Code, http.StatusInternalServerError)
}
