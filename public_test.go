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

var currencyResponseBody = `[
  {
    "currencyCodeA": 840,
    "currencyCodeB": 980,
    "date": 1552392228,
    "rateSell": 27,
    "rateBuy": 27.2,
    "rateCross": 27.1
  }
]`

var failResponseBody = `{
  "errorDescription": "go away"
}`

var expectedCurrencyResponseBody = []mono.CurrencyInfo{{
	CurrencyCodeAISO4217: 840,
	CurrencyCodeBISO4217: 980,
	Date:                 1552392228,
	RateSell:             27,
	RateBuy:              27.2,
	RateCross:            27.1,
}}

func TestNewPublic_Domain(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(currencyResponseBody))),
	}

	public := mono.NewPublic(mono.WithClient(client))

	actual, err := public.Currency(context.Background())
	if err != nil {
		t.Fatalf("No error expected, got: %v", err)
	}

	if !reflect.DeepEqual(actual, expectedCurrencyResponseBody) {
		t.Error("Actual != expected")
	}

	if client.Req.URL.Scheme != "https" {
		t.Error("Actual domain differs from expected. Actual: " + client.Req.URL.Scheme)
	}

	if client.Req.URL.Host != "api.monobank.ua" {
		t.Error("Actual domain differs from expected. Actual: " + client.Req.URL.Host)
	}
}

func TestPublic_Currency_Succ(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(currencyResponseBody))),
	}

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	actual, err := public.Currency(ctx)
	if err != nil {
		t.Fatalf("No error expected, got: %v", err)
	}

	if !reflect.DeepEqual(actual, expectedCurrencyResponseBody) {
		t.Error("Actual != expected")
	}
}

func TestPublic_Currency_FailMono(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(failResponseBody))),
	}

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if err.Error() != "mono error: go away" {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}
}
