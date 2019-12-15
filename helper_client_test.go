package mono_test

import (
	"net/http"
)

type currencyClient struct {
	Req  *http.Request
	Resp *http.Response
	Err  error
}

func (c *currencyClient) Do(req *http.Request) (*http.Response, error) {
	c.Req = req

	return c.Resp, c.Err
}
