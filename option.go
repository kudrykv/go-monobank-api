package mono

type optioner interface {
	setDomain(string)
	setClient(HTTPClient)
}

type Option func(optioner)

func WithDomain(domain string) Option {
	return func(o optioner) {
		o.setDomain(domain)
	}
}

func WithClient(client HTTPClient) Option {
	return func(o optioner) {
		o.setClient(client)
	}
}
