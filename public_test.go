package mono_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
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

var currencyFailResponseBody = `{
  "errorDescription": "go away"
}`

var expectedCurrencyResponseBody = []mono.CurrencyInfo{{
	CurrencyCodeA: 840,
	CurrencyCodeB: 980,
	Date:          1552392228,
	RateSell:      27,
	RateBuy:       27.2,
	RateCross:     27.1,
}}

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

	testCurrencyRequest(t, client.Req)
}

func TestPublic_Currency_FailMono(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(currencyFailResponseBody))),
	}

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if err.Error() != "mono error: go away" {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.Req)
}

func TestPublic_Currency_FailMalformedRequest(t *testing.T) {
	client := &currencyClient{}
	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain:invalid"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to create request: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}
}

func TestPublic_Currency_FailDoRequest(t *testing.T) {
	client := &currencyClient{}
	client.Err = errors.New("boo")

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to make request: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.Req)
}

func TestPublic_Currency_FailReadAll(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       &badReader{},
	}

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to read body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.Req)
}

func TestPublic_Currency_FailBodyClose(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       &badReadCloser{},
	}

	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to close the body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.Req)
}

func TestPublic_Currency_FailUnmarshalOnOK(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(currencyResponseBody))),
	}

	ush := &unmarshaller{}
	ush.Err = errors.New("boo")

	ctx := context.Background()

	public := mono.
		NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client), mono.WithUnmarshaller(ush))

	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to unmarshal body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.Req)
}

func TestPublic_Currency_FailUnmarshalOnBadRequest(t *testing.T) {
	client := &currencyClient{}
	client.Resp = &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(currencyFailResponseBody))),
	}

	ush := &unmarshaller{}
	ush.Err = errors.New("boo")

	ctx := context.Background()

	public := mono.
		NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client), mono.WithUnmarshaller(ush))

	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to unmarshal body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.Req)
}

func testCurrencyRequest(t *testing.T, req *http.Request) {
	if req.URL.Scheme != "https" {
		t.Error("Actual scheme differs from expected. Actual: " + req.URL.Scheme)
	}

	if req.URL.Host != "domain" {
		t.Error("Actual domain differs from expected. Actual: " + req.URL.Host)
	}

	if req.URL.Path != "/bank/currency" {
		t.Error("Actual scheme differs from expected. Actual: " + req.URL.Path)
	}
}
