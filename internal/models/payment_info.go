package models

import "time"

type PaymentInfo struct {
	ID         int64     `db:"id"          json:"id"`
	InvoiceID  int64     `db:"invoice_id"  json:"invoice_id"`
	CustomerID int64     `db:"customer_id" json:"customer_id"`
	BankName   string    `db:"bank_name"   json:"bank_name"`
	BankAccNo  string    `db:"bank_acc_no" json:"bank_acc_no"`
	BankBranch string    `db:"bank_branch" json:"bank_branch"`
	DueDate    time.Time `db:"due_date"    json:"due_date"`
	Notes      string    `db:"notes"       json:"notes,omitempty"`
}
