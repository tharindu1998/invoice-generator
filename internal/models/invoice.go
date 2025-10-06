package models

type Invoice struct {
    ID             int           `db:"id"`
    InvoiceNumber  string        `db:"invoice_number"`
    CustomerMobile string        `db:"customer_mobile"`
    CustomerName   string        `db:"customer_name"`
    CustomerEmail  string        `db:"customer_email"`
    Date           string        `db:"date"`
    DueDate        string        `db:"due_date"`
    TotalAmount    float64       `db:"total_amount"`
    Items          []InvoiceItem `db:"-"`
}
