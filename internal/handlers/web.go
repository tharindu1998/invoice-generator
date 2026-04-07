package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

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
		log.Printf("template parse error: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("template execute error: %v", err)
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
		log.Printf("get customers: %v", err)
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
		log.Printf("lookup customer: %v", err)
		jsonError(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// ── Save customer ────────────────────────────────────

func SaveCustomer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
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
		log.Printf("save customer: %v", err)
		http.Error(w, "could not save customer", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}

// ── Generate invoice ─────────────────────────────────

func GenerateInvoice(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
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
		qty, _ := strconv.Atoi(quantities[i])
		price, _ := strconv.ParseFloat(prices[i], 64)
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
		CustomerMobile: r.FormValue("customer_mobile"),
		CustomerName:   r.FormValue("customer_name"),
		CustomerEmail:  r.FormValue("customer_email"),
		Date:           date,
		PaymentDue:     dueDate,
		TotalAmount:    total,
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
		_ = repository.SaveCustomer(models.Customer{
			Phone: inv.CustomerMobile,
			Name:  inv.CustomerName,
			Email: inv.CustomerEmail,
		})
	}

	id, err := repository.CreateInvoice(inv, items, payment)
	if err != nil {
		log.Printf("create invoice: %v", err)
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
		log.Printf("get invoice: %v", err)
		http.NotFound(w, r)
		return
	}

	render(w, "invoice-preview.html", PageData{Page: "invoices", Data: inv})
}

// ── Invoices list ─────────────────────────────────────

func Invoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := repository.GetAllInvoices()
	if err != nil {
		log.Printf("get invoices: %v", err)
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
		log.Printf("get invoice for pdf: %v", err)
		http.NotFound(w, r)
		return
	}

	doc, err := pdf.GenerateInvoice(inv)
	if err != nil {
		log.Printf("generate pdf: %v", err)
		http.Error(w, "could not generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, inv.InvoiceNumber))
	if err := doc.Output(w); err != nil {
		log.Printf("pdf output: %v", err)
	}
}
