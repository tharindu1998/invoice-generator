package repository

import (
	"fmt"
	"time"

	"invoice-generator/internal/db"
	"invoice-generator/internal/models"
)

func CreateInvoice(inv models.Invoice, items []models.InvoiceItem, payment models.PaymentInfo) (int64, error) {
	tx, err := db.DB.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	inv.InvoiceNumber = generateInvoiceNumber()
	inv.Status = "draft"

	res, err := tx.NamedExec(`
		INSERT INTO invoices
			(invoice_number, customer_mobile, customer_name, customer_email, date, payment_due, total_amount, status)
		VALUES
			(:invoice_number, :customer_mobile, :customer_name, :customer_email, :date, :payment_due, :total_amount, :status)
	`, inv)
	if err != nil {
		return 0, err
	}

	invoiceID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		item.InvoiceID = invoiceID
		_, err = tx.NamedExec(`
			INSERT INTO invoice_items (invoice_id, name, quantity, price, amount)
			VALUES (:invoice_id, :name, :quantity, :price, :amount)
		`, item)
		if err != nil {
			return 0, err
		}
	}

	if payment.BankName != "" || payment.BankAccNo != "" || payment.Notes != "" {
		payment.InvoiceID = invoiceID
		_, err = tx.NamedExec(`
			INSERT INTO payment_info (invoice_id, bank_name, bank_acc_no, bank_branch, due_date, notes)
			VALUES (:invoice_id, :bank_name, :bank_acc_no, :bank_branch, :due_date, :notes)
		`, payment)
		if err != nil {
			return 0, err
		}
	}

	return invoiceID, tx.Commit()
}

func GetAllInvoices() ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := db.DB.Select(&invoices, "SELECT * FROM invoices ORDER BY created_at DESC")
	return invoices, err
}

func GetInvoice(id int64) (models.Invoice, error) {
	var inv models.Invoice
	if err := db.DB.Get(&inv, "SELECT * FROM invoices WHERE id = ?", id); err != nil {
		return inv, err
	}

	var items []models.InvoiceItem
	if err := db.DB.Select(&items, "SELECT * FROM invoice_items WHERE invoice_id = ?", id); err != nil {
		return inv, err
	}
	inv.Items = items

	var payment models.PaymentInfo
	err := db.DB.Get(&payment, "SELECT * FROM payment_info WHERE invoice_id = ? LIMIT 1", id)
	if err == nil {
		inv.Payment = &payment
	}

	return inv, nil
}

func DeleteInvoice(id int64) error {
	_, err := db.DB.Exec("DELETE FROM invoices WHERE id = ?", id)
	return err
}

func generateInvoiceNumber() string {
	var count int
	db.DB.Get(&count, "SELECT COUNT(*) FROM invoices WHERE DATE(created_at) = CURDATE()")
	return fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), count+1)
}
