package mono_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
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

	if p == nil {
		t.Fatal("personal must not be nil")
	}
}

func TestNewPersonal_PanicOnEmptyToken(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expected error, but got nil")
		}

		if err != "api token is required" {
			t.Error("got: " + err.(string))
		}
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
	if err != nil {
		t.Fatalf("No error expected, got: %v", err)
	}

	if !reflect.DeepEqual(actual, &expectedPersonal) {
		t.Error("Actual != expected")
	}
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
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if err.Error() != "mono error: go away" {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}
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
	if err != nil {
		t.Fatal("Expected err to be nil, got: " + err.Error())
	}

	if !reflect.DeepEqual(statements, expectedStatementsResponse) {
		t.Fatal("actual != expected")
	}
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

func expectError(t *testing.T, actual error, expected string) {
	if actual == nil {
		t.Error("Expected error to be present, actual nil")
		return
	}

	if actual.Error() != expected {
		t.Error("Actual error is '" + actual.Error() + "', expected '" + expected + "'")
	}
}
