package mono_test

import (
	"net/http"
)

type clienttest struct {
	Req  *http.Request
	Resp *http.Response
	Err  error
}

func (c *clienttest) Do(req *http.Request) (*http.Response, error) {
	c.Req = req

	return c.Resp, c.Err
}
