package mono

type optioner interface {
	setDomain(string)
	setClient(HTTPClient)
	setUnmarshaller(Unmarshaller)
	setWebhookBufferSize(uint32)
}

// Option allows to change default values for client.
type Option func(optioner)

// WithDomain allows to change the default API domain.
// The domain should be in format `scheme://domain`
func WithDomain(domain string) Option {
	return func(o optioner) {
		o.setDomain(domain)
	}
}

// WithClient allows to change default http.Client
func WithClient(client HTTPClient) Option {
	return func(o optioner) {
		o.setClient(client)
	}
}

// WithUnmarshaller allows to change default `json.Unmarshal` to something else.
func WithUnmarshaller(u Unmarshaller) Option {
	return func(o optioner) {
		o.setUnmarshaller(u)
	}
}

// WithWebhookBufferSize allows to change default buffer size.
// Default value is 100.
func WithWebhookBufferSize(size uint32) Option {
	return func(o optioner) {
		o.setWebhookBufferSize(size)
	}
}
