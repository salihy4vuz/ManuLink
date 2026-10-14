package services

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/go-pdf/fpdf"
)

// GenerateProformaPDF creates a professional B2B quotation PDF
// Returns the PDF as a byte slice so it can be sent via email or downloaded directly.
func GenerateProformaPDF(req PricingRequest, res PricingResponse, customerName string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// 1. Header (Company Logo & Info)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "SHEET METAL APEX INC.")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "123 Industrial Park, Manufacturing City")
	pdf.Ln(5)
	pdf.Cell(40, 10, "Email: sales@sheetmetalapex.com | Phone: +1-555-0198")
	pdf.Ln(15)

	// 2. Title & Date
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "OFFICIAL PROFORMA INVOICE (QUOTATION)")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", customerName))
	pdf.Ln(6)
	pdf.Cell(40, 10, fmt.Sprintf("Target Currency: %s", res.Currency))
	pdf.Ln(15)

	// 3. Technical Specifications
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Technical Specifications:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 8, fmt.Sprintf("- Material: %s (%.2f mm)", req.Material, req.ThicknessMm))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("- Quantity: %d units", req.Quantity))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("- Total Weight: %.2f kg (Includes %.1f%% scrap margin)", res.MaterialWeightKg, req.ScrapMarginPct))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("- Bending Operations: %d per unit", req.BendCount))
	pdf.Ln(15)

	// 4. Financial Breakdown Table
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Cost Breakdown:")
	pdf.Ln(10)
	
	// Table Headers
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(80, 8, "Description", "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, "Cost", "1", 0, "R", false, 0, "")
	pdf.Ln(-1)

	// Table Rows
	pdf.SetFont("Arial", "", 10)
	
	pdf.CellFormat(80, 8, "Material Cost", "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%.2f %s", res.MaterialCostUSD, res.Currency), "1", 0, "R", false, 0, "")
	pdf.Ln(-1)

	pdf.CellFormat(80, 8, "Laser Cutting Operations", "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%.2f %s", res.LaserCostUSD, res.Currency), "1", 0, "R", false, 0, "")
	pdf.Ln(-1)

	pdf.CellFormat(80, 8, "Bending Operations", "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%.2f %s", res.BendingCostUSD, res.Currency), "1", 0, "R", false, 0, "")
	pdf.Ln(-1)

	pdf.CellFormat(80, 8, "Machine Setup (Amortized)", "1", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%.2f %s", res.SetupCostUSD, res.Currency), "1", 0, "R", false, 0, "")
	pdf.Ln(10)

	// Final Price
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(80, 10, "TOTAL PROFORMA AMOUNT:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f %s", res.FinalSellingPrice, res.Currency), "0", 0, "R", false, 0, "")
	pdf.Ln(20)

	// 5. Footer & Signature
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(40, 10, "This is an automatically generated proforma invoice by SheetMetal Apex Engine.")
	pdf.Ln(5)
	pdf.Cell(40, 10, "Prices are valid for 24 hours due to LME material market volatility.")

	// Render to bytes buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		log.Printf("Failed to generate PDF: %v\n", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
