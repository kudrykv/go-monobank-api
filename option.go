package mono

type optioner interface {
	setDomain(domain string)
}

type Option func(optioner)
