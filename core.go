package mono

import "net/http"

type core struct {
	domain       string
	whBufferSize uint32

	tinyClient
}

func (c *core) setDomain(domain string) {
	c.domain = domain
}

func (c *core) setClient(client HTTPClient) {
	c.client = client
}

func (c *core) setUnmarshaller(u Unmarshaller) {
	c.unmarshaller = u
}

func (c *core) setWebhookBufferSize(size uint32) {
	c.whBufferSize = size
}

func newCore(opts ...Option) core {
	c := core{
		domain:       DefaultDomain,
		whBufferSize: 100,
		tinyClient: tinyClient{
			client:       &http.Client{},
			unmarshaller: unmarshaller{},
		},
	}

	for _, opt := range opts {
		opt(&c)
	}

	return c
}
