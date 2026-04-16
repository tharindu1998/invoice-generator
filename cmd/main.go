package main

import (
	"net/http"
	"os"

	"invoice-generator/internal/db"
	"invoice-generator/internal/handlers"
	"invoice-generator/internal/logger"
)

func main() {
	if _, err := db.InitDB(); err != nil {
		logger.Log.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	logger.Log.Info("database connected")

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Pages
	http.HandleFunc("GET /{$}", handlers.CreateInvoice)
	http.HandleFunc("GET /customers", handlers.Customers)
	http.HandleFunc("GET /invoices", handlers.Invoices)
	http.HandleFunc("GET /invoice/{id}", handlers.ViewInvoice)
	http.HandleFunc("GET /invoice/{id}/pdf", handlers.DownloadPDF)

	// Actions
	http.HandleFunc("POST /customers/save", handlers.SaveCustomer)
	http.HandleFunc("POST /customers/{id}/delete", handlers.DeleteCustomer)
	http.HandleFunc("POST /invoice/{id}/delete", handlers.DeleteInvoice)
	http.HandleFunc("POST /invoice/generate", handlers.GenerateInvoice)

	// API
	http.HandleFunc("GET /api/customer", handlers.LookupCustomer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.Info("server running", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Log.Error("server error", "err", err)
		os.Exit(1)
	}
}
