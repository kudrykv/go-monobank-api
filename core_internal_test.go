package mono

import (
	"net/http"
	"testing"
)

func TestCore_setDomain(t *testing.T) {
	c := newCore(WithDomain("test.com"))

	if c.domain != "test.com" {
		t.Fatal("Expected domain to be test.com, actual: " + c.domain)
	}
}

func TestCore_setClient(t *testing.T) {
	hct := &httpclienttest{}
	c := newCore(WithClient(hct))

	if c.client != hct {
		t.Fatal("expected client to be custom")
	}
}

func TestCore_setUnmarshaller(t *testing.T) {
	umhs := &unmtest{}
	c := newCore(WithUnmarshaller(umhs))

	if c.unmarshaller != umhs {
		t.Fatal("expected unmarshaller to be custom")
	}
}

func TestCore_setWebhookBufferSize(t *testing.T) {
	c := newCore(WithWebhookBufferSize(150))

	if c.whBufferSize != 150 {
		t.Fatal("expected wh buffer size to be custom")
	}
}

func TestCoreDefaults(t *testing.T) {
	c := newCore()

	if c.domain != DefaultDomain {
		t.Error("expected default domain, got: " + c.domain)
	}

	_, ok := c.client.(*http.Client)
	if !ok {
		t.Error("expected client to be a http.Client as default")
	}

	umsh := unmarshaller{}
	if c.unmarshaller != umsh {
		t.Error("expected default unmarshaller, got else")
	}

	if c.whBufferSize != 100 {
		t.Error("expected default wh buffer size, got else")
	}
}
