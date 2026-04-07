package pdf

import (
	"fmt"

	"github.com/go-pdf/fpdf"
	"invoice-generator/internal/models"
)

const (
	pageW    = 210.0
	marginL  = 15.0
	marginR  = 15.0
	contentW = pageW - marginL - marginR

	darkR, darkG, darkB       = 20, 20, 40
	mutedR, mutedG, mutedB    = 110, 110, 130
	borderR, borderG, borderB = 220, 220, 228
	accentR, accentG, accentB = 30, 150, 140
)

func line(pdf *fpdf.Fpdf, y float64) {
	pdf.SetDrawColor(borderR, borderG, borderB)
	pdf.SetLineWidth(0.25)
	pdf.Line(marginL, y, pageW-marginR, y)
}

func GenerateInvoice(inv models.Invoice) (*fpdf.Fpdf, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginL, marginL, marginR)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	// ── Title ────────────────────────────────────────────
	pdf.SetFont("Helvetica", "B", 26)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.SetXY(marginL, 12)
	pdf.CellFormat(contentW/2, 10, "Invoice", "", 2, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 8.5)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	pdf.CellFormat(contentW/2, 4, fmt.Sprintf("Invoice Number  #%s", inv.InvoiceNumber), "", 0, "L", false, 0, "")

	// ── Icon grid ────────────────────────────────────────
	ix, iy, sz, gap := pageW-marginR-17.0, 12.0, 7.5, 1.5
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			if row == 0 && col == 1 {
				pdf.SetFillColor(90, 190, 180)
			} else {
				pdf.SetFillColor(accentR, accentG, accentB)
			}
			pdf.RoundedRect(ix+float64(col)*(sz+gap), iy+float64(row)*(sz+gap), sz, sz, 1.8, "1234", "F")
		}
	}

	line(pdf, 30)

	// ── Billed By / Billed To ────────────────────────────
	half := contentW / 2
	rx := marginL + half

	// labels
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	pdf.SetXY(marginL, 34)
	pdf.CellFormat(half, 4, "Billed By :", "", 0, "L", false, 0, "")
	pdf.SetXY(rx, 34)
	pdf.CellFormat(half, 4, "Billed To :", "", 0, "L", false, 0, "")

	// seller
	sy := 40.0
	if inv.SellerName != "" {
		pdf.SetFont("Helvetica", "B", 12)
		pdf.SetTextColor(darkR, darkG, darkB)
		pdf.SetXY(marginL, sy)
		pdf.CellFormat(half-4, 6, inv.SellerName, "", 2, "L", false, 0, "")
		sy = pdf.GetY()
	}
	pdf.SetFont("Helvetica", "", 8.5)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	if inv.SellerPhone != "" {
		pdf.SetXY(marginL, sy)
		pdf.CellFormat(half-4, 4.5, inv.SellerPhone, "", 2, "L", false, 0, "")
		sy = pdf.GetY()
	}
	if inv.SellerAddress != "" {
		pdf.SetXY(marginL, sy)
		pdf.MultiCell(half-4, 4.5, inv.SellerAddress, "", "L", false)
		sy = pdf.GetY()
	}

	// customer
	cy := 40.0
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.SetXY(rx, cy)
	pdf.CellFormat(half, 6, inv.CustomerName, "", 2, "L", false, 0, "")
	cy = pdf.GetY()
	pdf.SetFont("Helvetica", "", 8.5)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	for _, v := range []string{inv.CustomerMobile, inv.CustomerEmail, inv.CustomerAddressLine1, inv.CustomerAddressLine2} {
		if v != "" {
			pdf.SetXY(rx, cy)
			pdf.CellFormat(half, 4.5, v, "", 2, "L", false, 0, "")
			cy = pdf.GetY()
		}
	}

	// dates
	dateY := sy
	if cy > dateY {
		dateY = cy
	}
	dateY += 4

	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	pdf.SetXY(marginL, dateY)
	pdf.CellFormat(half, 4, "Date Issued :", "", 0, "L", false, 0, "")
	pdf.SetXY(rx, dateY)
	pdf.CellFormat(half, 4, "Due Date:", "", 0, "L", false, 0, "")
	dateY += 5
	pdf.SetFont("Helvetica", "B", 10.5)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.SetXY(marginL, dateY)
	pdf.CellFormat(half, 6, inv.Date.Format("January 02, 2006"), "", 0, "L", false, 0, "")
	pdf.SetXY(rx, dateY)
	pdf.CellFormat(half, 6, inv.PaymentDue.Format("January 02, 2006"), "", 0, "L", false, 0, "")
	dateY += 8

	line(pdf, dateY)

	// ── Items table ──────────────────────────────────────
	y := dateY + 5

	pdf.SetFont("Helvetica", "", 8.5)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	pdf.SetXY(marginL, y)
	pdf.CellFormat(contentW, 5, "Invoice Details", "", 0, "L", false, 0, "")
	y += 7

	cDesc := 88.0
	cQty  := 22.0
	cUnit := 35.0
	cAmt  := contentW - cDesc - cQty - cUnit

	// column headers
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	pdf.SetXY(marginL, y)
	pdf.CellFormat(cDesc, 5, "Items/Service", "", 0, "L", false, 0, "")
	pdf.CellFormat(cQty,  5, "Quantity",      "", 0, "C", false, 0, "")
	pdf.CellFormat(cUnit, 5, "Unit Price",    "", 0, "R", false, 0, "")
	pdf.CellFormat(cAmt,  5, "Total",         "", 0, "R", false, 0, "")
	y += 5
	line(pdf, y)
	y += 2

	// rows
	for _, item := range inv.Items {
		pdf.SetFont("Helvetica", "", 9.5)
		pdf.SetTextColor(darkR, darkG, darkB)
		pdf.SetXY(marginL, y)
		pdf.CellFormat(cDesc, 6.5, item.Name, "", 0, "L", false, 0, "")
		pdf.CellFormat(cQty,  6.5, fmt.Sprintf("%d", item.Quantity), "", 0, "C", false, 0, "")
		pdf.CellFormat(cUnit, 6.5, fmt.Sprintf("%.2f", item.Price), "", 0, "R", false, 0, "")
		pdf.CellFormat(cAmt,  6.5, fmt.Sprintf("%.2f", item.Amount), "", 0, "R", false, 0, "")
		y += 6.5
		line(pdf, y)
		y += 2
	}

	// grand total
	y += 3
	pdf.SetFont("Helvetica", "B", 10.5)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.SetXY(marginL, y)
	pdf.CellFormat(cDesc+cQty+cUnit, 6, "Grand Total", "", 0, "R", false, 0, "")
	pdf.CellFormat(cAmt, 6, fmt.Sprintf("%.2f", inv.TotalAmount), "", 0, "R", false, 0, "")
	y += 9

	// ── Payment / Notes ──────────────────────────────────
	if inv.Payment != nil {
		p := inv.Payment
		if p.BankName != "" || p.BankAccNo != "" || p.BankBranch != "" {
			line(pdf, y)
			y += 5
			pdf.SetFont("Helvetica", "B", 8.5)
			pdf.SetTextColor(darkR, darkG, darkB)
			pdf.SetXY(marginL, y)
			pdf.CellFormat(contentW, 5, "Payment Details", "", 2, "L", false, 0, "")
			y += 5
			for _, pair := range [][2]string{{"Bank:", p.BankName}, {"Account No:", p.BankAccNo}, {"Branch:", p.BankBranch}} {
				if pair[1] != "" {
					pdf.SetFont("Helvetica", "", 8.5)
					pdf.SetTextColor(mutedR, mutedG, mutedB)
					pdf.SetXY(marginL, y)
					pdf.CellFormat(35, 5, pair[0], "", 0, "L", false, 0, "")
					pdf.SetTextColor(darkR, darkG, darkB)
					pdf.CellFormat(contentW-35, 5, pair[1], "", 0, "L", false, 0, "")
					y += 5
				}
			}
			y += 3
		}
		if p.Notes != "" {
			boxH := 18.0
			pdf.SetFillColor(248, 249, 250)
			pdf.SetDrawColor(borderR, borderG, borderB)
			pdf.RoundedRect(marginL, y, contentW, boxH, 2.5, "1234", "FD")
			pdf.SetFont("Helvetica", "B", 8)
			pdf.SetTextColor(darkR, darkG, darkB)
			pdf.SetXY(marginL+3, y+3)
			pdf.CellFormat(contentW-6, 4, "Notes", "", 2, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 7.5)
			pdf.SetTextColor(mutedR, mutedG, mutedB)
			pdf.SetXY(marginL+3, y+8)
			pdf.MultiCell(contentW-6, 4, p.Notes, "", "L", false)
		}
	}

	// ── Footer ───────────────────────────────────────────
	footerY := 282.0
	line(pdf, footerY)
	pdf.SetFont("Helvetica", "B", 8.5)
	pdf.SetTextColor(darkR, darkG, darkB)
	pdf.SetXY(marginL, footerY+3)
	pdf.CellFormat(contentW/2, 5, inv.SellerName, "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8.5)
	pdf.SetTextColor(mutedR, mutedG, mutedB)
	pdf.CellFormat(contentW/2, 5, inv.SellerPhone, "", 0, "R", false, 0, "")

	return pdf, pdf.Error()
}
