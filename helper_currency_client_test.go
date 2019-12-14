package mono_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type succCurrencyClient struct {
	req *http.Request
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
