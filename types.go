package mono

type Cashback string

const (
	CashbackNone  Cashback = "None"
	CashbackUAH   Cashback = "UAH"
	CashbackMiles Cashback = "Miles"
)

type UserInfo struct {
	Name       string    `json:"name"`
	WebHookURL string    `json:"webHookUrl"`
	Accounts   []Account `json:"accounts"`
}

type Account struct {
	ID           string   `json:"id"`
	Balance      int64    `json:"balance"`
	CreditLimit  int64    `json:"creditLimit"`
	CurrencyCode int      `json:"currencyCode"`
	CashbackType Cashback `json:"cashbackType"`
}

type StatementItem struct {
	ID          string `json:"id"`
	Time        int64  `json:"time"`
	Description string `json:"description"`
	// Merchant Category Code
	MCC             int   `json:"mcc"`
	Hold            bool  `json:"hold"`
	Amount          int64 `json:"amount"`
	OperationAmount int64 `json:"operationAmount"`
	CurrencyCode    int   `json:"currencyCode"`
	CommissionRate  int64 `json:"commissionRate"`
	CashbackAmount  int64 `json:"cashbackAmount"`
	Balance         int64 `json:"balance"`
}

type CurrencyInfo struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	Date          int64   `json:"date"`
	RateSell      float64 `json:"rateSell"`
	RateBuy       float64 `json:"rateBuy"`
	RateCross     float64 `json:"rateCross"`
}

type ErrorMono struct {
	Description string `json:"errorDescription"`
}
