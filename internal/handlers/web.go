package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"invoice-generator/internal/logger"
	"invoice-generator/internal/models"
	"invoice-generator/internal/pdf"
	"invoice-generator/internal/repository"
)

type PageData struct {
	Page string
	Data any
}

// ── Template rendering ───────────────────────────────

func render(w http.ResponseWriter, page string, data PageData) {
	tmpl, err := template.ParseFiles(
		filepath.Join("web", "templates", "base.html"),
		filepath.Join("web", "templates", page),
	)
	if err != nil {
		logger.Log.Error("template parse error", "err", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		logger.Log.Error("template execute error", "err", err)
	}
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// ── Page handlers ────────────────────────────────────

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	render(w, "create-invoice.html", PageData{Page: "invoice"})
}

func Customers(w http.ResponseWriter, r *http.Request) {
	customers, err := repository.GetAllCustomers()
	if err != nil {
		logger.Log.Error("get customers", "err", err)
		customers = nil
	}
	render(w, "customer-form.html", PageData{Page: "customers", Data: customers})
}

// ── API: customer lookup by phone ────────────────────

func LookupCustomer(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		jsonError(w, "phone required", http.StatusBadRequest)
		return
	}

	c, err := repository.GetCustomerByPhone(phone)
	if errors.Is(err, repository.ErrNotFound) {
		jsonError(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		logger.Log.Error("lookup customer", "err", err)
		jsonError(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// ── Save customer ────────────────────────────────────

func SaveCustomer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Log.Error("parse form", "err", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	c := models.Customer{
		Phone:        r.FormValue("phone"),
		Name:         r.FormValue("name"),
		Email:        r.FormValue("email"),
		AddressLine1: r.FormValue("address_line1"),
		AddressLine2: r.FormValue("address_line2"),
	}

	if c.Phone == "" || c.Name == "" {
		http.Error(w, "phone and name are required", http.StatusBadRequest)
		return
	}

	if err := repository.SaveCustomer(c); err != nil {
		logger.Log.Error("save customer", "err", err)
		http.Error(w, "could not save customer", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}

// ── Delete customer ──────────────────────────────────

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteCustomer(id); err != nil {
		logger.Log.Error("delete customer", "err", err)
		http.Error(w, "could not delete customer", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}

// ── Delete invoice ───────────────────────────────────

func DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := repository.DeleteInvoice(id); err != nil {
		logger.Log.Error("delete invoice", "err", err)
		http.Error(w, "could not delete invoice", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/invoices", http.StatusSeeOther)
}

// ── Generate invoice ─────────────────────────────────

func GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger.Log.Error("parse form", "err", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}
	dueDate, err := time.Parse("2006-01-02", r.FormValue("due_date"))
	if err != nil {
		http.Error(w, "invalid due date", http.StatusBadRequest)
		return
	}

	names := r.Form["item_name[]"]
	quantities := r.Form["item_quantity[]"]
	prices := r.Form["item_price[]"]

	var items []models.InvoiceItem
	var total float64

	for i := range names {
		qty, err := strconv.Atoi(quantities[i])
		if err != nil {
			logger.Log.Error("parse quantity", "err", err)
			http.Error(w, "invalid quantity", http.StatusBadRequest)
			return
		}
		price, err := strconv.ParseFloat(prices[i], 64)
		if err != nil {
			logger.Log.Error("parse price", "err", err)
			http.Error(w, "invalid price", http.StatusBadRequest)
			return
		}
		amount := float64(qty) * price
		total += amount

		items = append(items, models.InvoiceItem{
			Name:     names[i],
			Quantity: qty,
			Price:    price,
			Amount:   amount,
		})
	}

	inv := models.Invoice{
		CustomerMobile:       r.FormValue("customer_mobile"),
		CustomerName:         r.FormValue("customer_name"),
		CustomerEmail:        r.FormValue("customer_email"),
		CustomerAddressLine1: r.FormValue("customer_address_line1"),
		CustomerAddressLine2: r.FormValue("customer_address_line2"),
		SellerName:           r.FormValue("seller_name"),
		SellerPhone:          r.FormValue("seller_phone"),
		SellerAddress:        r.FormValue("seller_address"),
		Date:                 date,
		PaymentDue:           dueDate,
		TotalAmount:          total,
	}

	payment := models.PaymentInfo{
		BankName:   r.FormValue("bank_name"),
		BankAccNo:  r.FormValue("bank_acc_no"),
		BankBranch: r.FormValue("bank_branch"),
		DueDate:    dueDate,
		Notes:      r.FormValue("notes"),
	}

	// Auto-save customer by phone if name provided
	if inv.CustomerMobile != "" && inv.CustomerName != "" {
		if err := repository.SaveCustomer(models.Customer{
			Phone:        inv.CustomerMobile,
			Name:         inv.CustomerName,
			Email:        inv.CustomerEmail,
			AddressLine1: inv.CustomerAddressLine1,
			AddressLine2: inv.CustomerAddressLine2,
		}); err != nil {
			logger.Log.Error("auto-save customer", "err", err)
		}
	}

	id, err := repository.CreateInvoice(inv, items, payment)
	if err != nil {
		logger.Log.Error("create invoice", "err", err)
		http.Error(w, "could not save invoice", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/invoice/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

// ── View invoice ─────────────────────────────────────

func ViewInvoice(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	inv, err := repository.GetInvoice(id)
	if err != nil {
		logger.Log.Error("get invoice", "err", err)
		http.NotFound(w, r)
		return
	}

	render(w, "invoice-preview.html", PageData{Page: "invoices", Data: inv})
}

// ── Invoices list ─────────────────────────────────────

func Invoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := repository.GetAllInvoices()
	if err != nil {
		logger.Log.Error("get invoices", "err", err)
		invoices = nil
	}
	render(w, "invoices.html", PageData{Page: "invoices", Data: invoices})
}

// ── Download PDF ──────────────────────────────────────

func DownloadPDF(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	inv, err := repository.GetInvoice(id)
	if err != nil {
		logger.Log.Error("get invoice for pdf", "err", err)
		http.NotFound(w, r)
		return
	}

	doc, err := pdf.GenerateInvoice(inv)
	if err != nil {
		logger.Log.Error("generate pdf", "err", err)
		http.Error(w, "could not generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, inv.InvoiceNumber))
	if err := doc.Output(w); err != nil {
		logger.Log.Error("pdf output", "err", err)
	}
}
