package mono

import (
	"time"
)

// Time wraps int64 to add a func for converting the value to time.Time.
type Time int64

// Time creates the time.Time from the int64 value.
func (t Time) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// UserInfo describes customer and customer's accounts.
type UserInfo struct {
	// Name describes client name.
	Name string `json:"name"`
	// WebHookURL for getting information about the new transaction.
	WebHookURL string `json:"webHookUrl"`
	// Accounts list available accounts.
	Accounts []Account `json:"accounts"`
}

// Account describes customer's account.
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
	Time        Time   `json:"time"`
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
	// Commission rate in transaction's currency in the minimal units -- cents of the corresponding currency.
	CommissionRate int64 `json:"commissionRate"`
	// Cashback amount in account currency in the minimal units -- cents of the corresponding currency.
	CashbackAmount int64 `json:"cashbackAmount"`
	// Balance in the minimal units -- cents of the corresponding currency.
	Balance int64 `json:"balance"`
}

// CurrencyInfo specifies single currency rate.
type CurrencyInfo struct {
	CurrencyCodeAISO4217 int `json:"currencyCodeA"`
	CurrencyCodeBISO4217 int `json:"currencyCodeB"`
	// Rate at the given point of time in UNIX timestamp.
	Date      Time    `json:"date"`
	RateSell  float64 `json:"rateSell"`
	RateBuy   float64 `json:"rateBuy"`
	RateCross float64 `json:"rateCross"`
}

// WebhookData defines the shape of the incoming webhook object.
type WebhookData struct {
	Type string `json:"type"`

	Data WebhookStatementItem `json:"data"`
}

// WebhookStatementItem is the transaction item from the webhook.
type WebhookStatementItem struct {
	AccountID     string        `json:"account"`
	StatementItem StatementItem `json:"statementItem"`
}

type errorMono struct {
	Description string `json:"errorDescription"`
}
