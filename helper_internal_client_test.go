package mono

import (
	"errors"
	"io"
	"net/http"
)

type httpclienttest struct {
	Req  *http.Request
	Resp *http.Response
	Err  error
}

func (c *httpclienttest) Do(req *http.Request) (*http.Response, error) {
	c.Req = req

	return c.Resp, c.Err
}

type badReader struct {
}

func (b badReader) Read([]byte) (int, error) {
	return 0, errors.New("boo")
}

func (b badReader) Close() error {
	panic("reading failed already")
}

type badReadCloser struct {
}

func (b badReadCloser) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (b badReadCloser) Close() error {
	return errors.New("boo")
}

type unmtest struct {
	Err error
}

func (b unmtest) Unmarshal([]byte, interface{}) error {
	return b.Err
}
