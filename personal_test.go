package mono_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

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
	client := &currencyClient{}
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
	client := &currencyClient{}
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
