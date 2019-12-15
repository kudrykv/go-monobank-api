package mono_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
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
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(currencyResponseBody))),
	}

	public := mono.NewPublic(mono.WithClient(client))

	actual, err := public.Currency(context.Background())
	expectNoError(t, err)
	expectDeepEquals(t, actual, expectedCurrencyResponseBody)
	expectEquals(t, client.Req.URL.Scheme, "https")
	expectEquals(t, client.Req.URL.Host, "api.monobank.ua")
}

func TestPublic_Currency_Succ(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(currencyResponseBody))),
	}

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	actual, err := public.Currency(ctx)
	expectNoError(t, err)
	expectDeepEquals(t, actual, expectedCurrencyResponseBody)
}

func TestPublic_Currency_FailMono(t *testing.T) {
	client := &clienttest{}
	client.Resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(failResponseBody))),
	}

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	expectError(t, err, "mono error: go away")
}
