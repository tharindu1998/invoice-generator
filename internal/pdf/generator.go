package pdf

import (
	"fmt"

	"github.com/go-pdf/fpdf"
	"invoice-generator/internal/models"
)

const (
	pageW      = 210.0
	marginL    = 15.0
	marginR    = 15.0
	contentW   = pageW - marginL - marginR
	primaryR   = 79
	primaryG   = 70
	primaryB   = 229
)

func GenerateInvoice(inv models.Invoice) (*fpdf.Fpdf, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginL, 15, marginR)
	pdf.AddPage()

	// ── Header bar ──────────────────────────────────────
	pdf.SetFillColor(primaryR, primaryG, primaryB)
	pdf.Rect(0, 0, pageW, 28, "F")

	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(marginL, 8)
	pdf.CellFormat(contentW/2, 12, "INVOICE", "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(marginL+contentW/2, 8)
	pdf.CellFormat(contentW/2, 6, fmt.Sprintf("# %s", inv.InvoiceNumber), "", 2, "R", false, 0, "")
	pdf.SetXY(marginL+contentW/2, 14)
	pdf.CellFormat(contentW/2, 6, fmt.Sprintf("Date: %s", inv.Date.Format("02 Jan 2006")), "", 0, "R", false, 0, "")

	// ── Seller (left) & Customer (right) ────────────────
	pdf.SetXY(marginL, 36)

	// Seller — left column
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(100, 100, 120)
	pdf.CellFormat(contentW/2, 5, "FROM", "", 2, "L", false, 0, "")
	if inv.SellerName != "" {
		pdf.SetFont("Helvetica", "B", 11)
		pdf.SetTextColor(30, 30, 30)
		pdf.CellFormat(contentW/2, 6, inv.SellerName, "", 2, "L", false, 0, "")
	}
	if inv.SellerAddress != "" {
		pdf.SetFont("Helvetica", "", 9)
		pdf.SetTextColor(80, 80, 90)
		pdf.MultiCell(contentW/2, 5, inv.SellerAddress, "", "L", false)
	}

	// Customer — right column
	pdf.SetXY(marginL+contentW/2, 36)
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(100, 100, 120)
	pdf.CellFormat(contentW/2, 5, "BILLED TO", "", 2, "R", false, 0, "")

	pdf.SetXY(marginL+contentW/2, 41)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(contentW/2, 6, inv.CustomerName, "", 2, "R", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(80, 80, 90)
	if inv.CustomerMobile != "" {
		pdf.CellFormat(contentW/2, 5, inv.CustomerMobile, "", 2, "R", false, 0, "")
	}
	if inv.CustomerEmail != "" {
		pdf.CellFormat(contentW/2, 5, inv.CustomerEmail, "", 2, "R", false, 0, "")
	}
	if inv.CustomerAddress != "" {
		pdf.CellFormat(contentW/2, 5, inv.CustomerAddress, "", 2, "R", false, 0, "")
	}

	// Due date — below right column
	pdf.SetXY(marginL+contentW/2, pdf.GetY()+3)
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(100, 100, 120)
	pdf.CellFormat(contentW/2, 5, "PAYMENT DUE", "", 2, "R", false, 0, "")
	pdf.SetXY(marginL+contentW/2, pdf.GetY())
	pdf.SetFont("Helvetica", "B", 13)
	pdf.SetTextColor(primaryR, primaryG, primaryB)
	pdf.CellFormat(contentW/2, 7, inv.PaymentDue.Format("02 Jan 2006"), "", 0, "R", false, 0, "")

	// ── Items table ──────────────────────────────────────
	y := 72.0
	pdf.SetXY(marginL, y)

	// Table header
	pdf.SetFillColor(245, 246, 250)
	pdf.SetDrawColor(220, 220, 230)
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(100, 100, 120)

	colW := [4]float64{85, 20, 32, 33}
	headers := [4]string{"DESCRIPTION", "QTY", "UNIT PRICE", "AMOUNT"}
	aligns := [4]string{"L", "C", "R", "R"}

	for i, h := range headers {
		pdf.CellFormat(colW[i], 8, h, "B", 0, aligns[i], true, 0, "")
	}
	pdf.Ln(-1)

	// Rows
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(30, 30, 30)
	fill := false
	for _, item := range inv.Items {
		if fill {
			pdf.SetFillColor(250, 250, 255)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.CellFormat(colW[0], 8, item.Name, "", 0, "L", true, 0, "")
		pdf.CellFormat(colW[1], 8, fmt.Sprintf("%d", item.Quantity), "", 0, "C", true, 0, "")
		pdf.CellFormat(colW[2], 8, fmt.Sprintf("%.2f", item.Price), "", 0, "R", true, 0, "")
		pdf.CellFormat(colW[3], 8, fmt.Sprintf("%.2f", item.Amount), "", 0, "R", true, 0, "")
		pdf.Ln(-1)
		fill = !fill
	}

	// Total row
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(primaryR, primaryG, primaryB)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(colW[0]+colW[1]+colW[2], 9, "TOTAL", "", 0, "R", true, 0, "")
	pdf.CellFormat(colW[3], 9, fmt.Sprintf("%.2f", inv.TotalAmount), "", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// ── Payment info ─────────────────────────────────────
	if inv.Payment != nil && (inv.Payment.BankName != "" || inv.Payment.BankAccNo != "") {
		p := inv.Payment
		pdf.Ln(8)
		pdf.SetFont("Helvetica", "B", 8)
		pdf.SetTextColor(100, 100, 120)
		pdf.SetFillColor(245, 246, 250)
		pdf.CellFormat(contentW, 6, "PAYMENT DETAILS", "T", 2, "L", false, 0, "")

		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(30, 30, 30)
		if p.BankName != "" {
			pdf.CellFormat(40, 6, "Bank:", "", 0, "L", false, 0, "")
			pdf.CellFormat(contentW-40, 6, p.BankName, "", 2, "L", false, 0, "")
		}
		if p.BankAccNo != "" {
			pdf.CellFormat(40, 6, "Account No:", "", 0, "L", false, 0, "")
			pdf.CellFormat(contentW-40, 6, p.BankAccNo, "", 2, "L", false, 0, "")
		}
		if p.BankBranch != "" {
			pdf.CellFormat(40, 6, "Branch:", "", 0, "L", false, 0, "")
			pdf.CellFormat(contentW-40, 6, p.BankBranch, "", 2, "L", false, 0, "")
		}
		if p.Notes != "" {
			pdf.Ln(3)
			pdf.SetFont("Helvetica", "I", 9)
			pdf.SetTextColor(100, 100, 120)
			pdf.MultiCell(contentW, 5, p.Notes, "", "L", false)
		}
	}

	// ── Footer ───────────────────────────────────────────
	pdf.SetY(-20)
	pdf.SetFont("Helvetica", "I", 8)
	pdf.SetTextColor(160, 160, 170)
	pdf.CellFormat(contentW, 5, "Generated by InvoiceGen", "", 0, "C", false, 0, "")

	return pdf, pdf.Error()
}
