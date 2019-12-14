package mono

type optioner interface {
	setDomain(string)
	setClient(HTTPClient)
}

type Option func(optioner)
