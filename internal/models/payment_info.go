package models

type PaymentInfo struct {
	InvoiceID   string `json:"invoice_id"`
	CustomerID  string `json:"customer_id"`
	BankName    string `json:"bank_name"`
	BankAccNo   string `json:"bank_acc_no"`
	BankBranch  string `json:"bank_branch"`
	DueDate     string `json:"due_date"`
	Notes       string `json:"notes,omitempty"`
}
