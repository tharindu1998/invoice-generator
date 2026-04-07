package models

import "time"

type Invoice struct {
	ID             int64         `db:"id"              json:"id"`
	InvoiceNumber  string        `db:"invoice_number"  json:"invoice_number"`
	CustomerMobile string        `db:"customer_mobile" json:"customer_mobile"`
	CustomerName   string        `db:"customer_name"   json:"customer_name"`
	CustomerEmail  string        `db:"customer_email"  json:"customer_email"`
	Date           time.Time     `db:"date"            json:"date"`
	PaymentDue     time.Time     `db:"payment_due"     json:"payment_due"`
	TotalAmount    float64       `db:"total_amount"    json:"total_amount"`
	Status         string        `db:"status"          json:"status"`
	CreatedAt      time.Time     `db:"created_at"      json:"created_at"`
	Items          []InvoiceItem `db:"-"               json:"items"`
	Payment        *PaymentInfo  `db:"-"               json:"payment,omitempty"`
}
