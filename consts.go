package mono

// Cashback is the enum of allowed cashback types.
type Cashback string

const (
	// DefaultDomain is the domain used by default to call the bank.
	DefaultDomain = "https://api.monobank.ua"

	// CashbackNone tells there is no cashback.
	CashbackNone Cashback = "None"
	// CashbackUAH tells the cashback is in UAH.
	CashbackUAH Cashback = "UAH"
	// CashbackMiles tells the cashback is in Miles.
	CashbackMiles Cashback = "Miles"
)
