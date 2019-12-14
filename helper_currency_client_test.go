package mono_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

type succCurrencyClient struct {
	basicClient
}

func (c *succCurrencyClient) Do(req *http.Request) (*http.Response, error) {
	c.req = req

	body := `[
  {
    "currencyCodeA": 840,
    "currencyCodeB": 980,
    "date": 1552392228,
    "rateSell": 27,
    "rateBuy": 27.2,
    "rateCross": 27.1
  }
]`

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}, nil
}

func (c *succCurrencyClient) request() *http.Request {
	return c.req
}

type failCurrencyClient struct {
	basicClient
}

func (f *failCurrencyClient) Do(req *http.Request) (*http.Response, error) {
	f.req = req

	body := `{
  "errorDescription": "go away"
}`

	return &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
	}, nil
}

type doFailCurrencyClient struct {
	basicClient
}

func (d *doFailCurrencyClient) Do(req *http.Request) (*http.Response, error) {
	d.req = req

	return nil, errors.New("boo")
}

type basicClient struct {
	req *http.Request
}

func (c *basicClient) request() *http.Request {
	return c.req
}
