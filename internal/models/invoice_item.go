package models

type InvoiceItem struct {
    ID        int     `db:"id"`
    InvoiceID int     `db:"invoice_id"`
    ProductID int     `db:"product_id"`
    Quantity  int     `db:"quantity"`
    UnitPrice float64 `db:"unit_price"`
    Total     float64 `db:"total"`
}