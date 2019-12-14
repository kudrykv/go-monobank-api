package mono_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	mono "github.com/kudrykv/go-monobank-api"
)

func TestPublic_Currency_Succ(t *testing.T) {
	client := &succCurrencyClient{}
	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	actual, err := public.Currency(ctx)
	if err != nil {
		t.Fatalf("No error expected, got: %v", err)
	}

	expected := []mono.CurrencyInfo{{
		CurrencyCodeA: 840,
		CurrencyCodeB: 980,
		Date:          1552392228,
		RateSell:      27,
		RateBuy:       27.2,
		RateCross:     27.1,
	}}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("Actual != expected")
	}

	testCurrencyRequest(t, client.request())
}

func TestPublic_Currency_FailMono(t *testing.T) {
	client := &failCurrencyClient{}
	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if err.Error() != "mono error: go away" {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.request())
}

func TestPublic_Currency_FailMalformedRequest(t *testing.T) {
	client := &failCurrencyClient{}
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
	client := &doFailCurrencyClient{}
	ctx := context.Background()

	public := mono.NewPublic(mono.WithDomain("https://domain"), mono.WithClient(client))
	_, err := public.Currency(ctx)
	if err == nil {
		t.Fatal("No error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to make request: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testCurrencyRequest(t, client.request())
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
