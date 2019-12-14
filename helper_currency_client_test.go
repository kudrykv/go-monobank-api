package mono_test

import (
	"errors"
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

type badReader struct {
}

func (b badReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("boo")
}

func (b badReader) Close() error {
	return errors.New("boo")
}
