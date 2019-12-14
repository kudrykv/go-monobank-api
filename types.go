package mono

type Cashback string

const (
	CashbackNone  Cashback = "None"
	CashbackUAH   Cashback = "UAH"
	CashbackMiles Cashback = "Miles"
)

// UserInfo describes customer and customer's accounts.
type UserInfo struct {
	// Name describes client name.
	Name string `json:"name"`
	// WebHookURL for getting information about the new transaction.
	WebHookURL string `json:"webHookUrl"`
	// Accounts list available accounts.
	Accounts []Account `json:"accounts"`
}

type Account struct {
	// Identifier of the account.
	ID string `json:"id"`
	// Balance in the minimal units -- cents of the corresponding currency.
	Balance int64 `json:"balance"`
	// Credit limit.
	CreditLimit int64 `json:"creditLimit"`
	// Currency code in ISO 4217.
	CurrencyCodeISO4217 int `json:"currencyCode"`
	// Type of the cashback.
	// Available values are `None`, `UAH`, and `Miles`.
	// One can refer using package's consts.
	CashbackType Cashback `json:"cashbackType"`
}

// StatementItem is the transaction entry.
type StatementItem struct {
	// Transaction identifier.
	ID string `json:"id"`
	// Time when the transaction was made in UNIX timestamp.
	Time        int64  `json:"time"`
	Description string `json:"description"`
	// Merchant Category Code
	MCC int `json:"mcc"`
	// Hold state.
	// Learn more: https://en.wikipedia.org/wiki/Authorization_hold
	Hold bool `json:"hold"`
	// Amount in account currency in the minimal units -- cents of the corresponding currency.
	Amount int64 `json:"amount"`
	// Amount in transaction currency in the minimal units -- cents of the corresponding currency.
	OperationAmount     int64 `json:"operationAmount"`
	CurrencyCodeISO4217 int   `json:"currencyCode"`
	// Comission rate in transaction's currency in the minimal units -- cents of the corresponding currency.
	CommissionRate int64 `json:"commissionRate"`
	// Cashback amount in account currency in the minimal units -- cents of the corresponding currency.
	CashbackAmount int64 `json:"cashbackAmount"`
	// Balance in the minimal units -- cents of the corresponding currency.
	Balance int64 `json:"balance"`
}

type CurrencyInfo struct {
	CurrencyCodeAISO4217 int `json:"currencyCodeA"`
	CurrencyCodeBISO4217 int `json:"currencyCodeB"`
	// Rate at the given point of time in UNIX timestamp.
	Date      int64   `json:"date"`
	RateSell  float64 `json:"rateSell"`
	RateBuy   float64 `json:"rateBuy"`
	RateCross float64 `json:"rateCross"`
}

type errorMono struct {
	Description string `json:"errorDescription"`
}
