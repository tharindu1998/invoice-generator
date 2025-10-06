package models

type InvoiceItem struct {
    ID        int64   `json:"id"`
    InvoiceID int64   `json:"invoice_id"`
    ProductID int64   `json:"product_id"`
    Name      string  `json:"name"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
    Amount    float64 `json:"amount"`
}