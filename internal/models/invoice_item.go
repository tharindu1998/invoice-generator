package models

type InvoiceItem struct {
	ID        int64   `db:"id"         json:"id"`
	InvoiceID int64   `db:"invoice_id" json:"invoice_id"`
	Name      string  `db:"name"       json:"name"`
	Quantity  int     `db:"quantity"   json:"quantity"`
	Price     float64 `db:"price"      json:"price"`
	Amount    float64 `db:"amount"     json:"amount"`
}
