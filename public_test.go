package mono_test

import (
	"context"
	"reflect"
	"testing"

	mono "github.com/kudrykv/go-monobank-api"
)

func TestPublic_Currency(t *testing.T) {
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
		t.Fatalf("Actual != expected")
	}

	req := client.request()

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
