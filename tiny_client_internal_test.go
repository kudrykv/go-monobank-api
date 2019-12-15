package mono

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const testFailResponseBody = `{"errorDescription":"go away"}`

func TestTinyClientRequest_Succ(t *testing.T) {
	hct := &httpclienttest{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`42`))),
		},
		Err: nil,
	}

	client := tinyClient{client: hct, unmarshaller: unmarshaller{}}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain/url", nil, &ultimateAnswer)
	if err != nil {
		t.Fatalf("No error expected, got: %v", err)
	}

	if ultimateAnswer != 42 {
		t.Fatalf("Expecte 42, got: %d", ultimateAnswer)
	}

	testRequest(t, hct.Req)
}

func TestTinyClientRequest_FailMono(t *testing.T) {
	hct := &httpclienttest{
		Resp: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(testFailResponseBody))),
		},
		Err: nil,
	}

	client := tinyClient{client: hct, unmarshaller: unmarshaller{}}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain/url", nil, &ultimateAnswer)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if err.Error() != "mono error: go away" {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testRequest(t, hct.Req)
}

func TestTinyClientRequest_FailMalformedRequest(t *testing.T) {
	hct := &httpclienttest{}
	client := tinyClient{client: hct, unmarshaller: unmarshaller{}}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain:err/url", nil, &ultimateAnswer)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to create request: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}
}

func TestTinyClientRequest_FailDoRequest(t *testing.T) {
	hct := &httpclienttest{Err: errors.New("boo")}
	client := tinyClient{client: hct, unmarshaller: unmarshaller{}}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain/url", nil, &ultimateAnswer)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to make request: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testRequest(t, hct.Req)
}

func TestTinyClientRequest_FailReadAll(t *testing.T) {
	hct := &httpclienttest{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &badReader{},
		},
		Err: nil,
	}

	client := tinyClient{client: hct, unmarshaller: unmarshaller{}}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain/url", nil, &ultimateAnswer)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to read body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testRequest(t, hct.Req)
}

func TestTinyClientRequest_FailBodyClose(t *testing.T) {
	hct := &httpclienttest{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &badReadCloser{},
		},
		Err: nil,
	}

	client := tinyClient{client: hct, unmarshaller: unmarshaller{}}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain/url", nil, &ultimateAnswer)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to close the body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testRequest(t, hct.Req)
}

func TestTinyClientRequest_FailUnmarshalOnBadRequest(t *testing.T) {
	hct := &httpclienttest{
		Resp: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(testFailResponseBody))),
		},
		Err: nil,
	}

	u := unmtest{}
	u.Err = errors.New("boo")

	client := tinyClient{client: hct, unmarshaller: u}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain/url", nil, &ultimateAnswer)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to unmarshal body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testRequest(t, hct.Req)
}

func TestTinyClientRequest_FailUnmarshalOnOK(t *testing.T) {
	hct := &httpclienttest{
		Resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`42`))),
		},
		Err: nil,
	}

	u := unmtest{}
	u.Err = errors.New("boo")

	client := tinyClient{client: hct, unmarshaller: u}

	var ultimateAnswer int
	err := client.request(context.Background(), http.MethodGet, "https://domain/url", nil, &ultimateAnswer)
	if err == nil {
		t.Fatal("Error expected, got nil")
	}

	if strings.Index(err.Error(), "failed to unmarshal body: ") != 0 {
		t.Error("Actual error differs from expected. Actual> " + err.Error())
	}

	testRequest(t, hct.Req)
}

func testRequest(t *testing.T, req *http.Request) {
	if req.URL.Scheme != "https" {
		t.Error("Actual scheme differs from expected. Actual: " + req.URL.Scheme)
	}

	if req.URL.Host != "domain" {
		t.Error("Actual domain differs from expected. Actual: " + req.URL.Host)
	}

	if req.URL.Path != "/url" {
		t.Error("Actual scheme differs from expected. Actual: " + req.URL.Path)
	}
}
