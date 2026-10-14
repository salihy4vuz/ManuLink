package services

import (
	"fmt"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/models"
)

// GenerateQuotePDF creates an A4 PDF for a given quote using the Tenant's branding
func GenerateQuotePDF(tenant models.Tenant, order models.Order) ([]byte, error) {
	m := maroto.New()

	// Use primary color (convert hex to RGB normally, hardcoded for simplicity)
	// You can add logic to parse tenant.Theme.PrimaryColor

	// Header Row
	m.AddRow(20,
		col.New(12).Add(
			text.New(tenant.Name+" - Official Quote", props.Text{
				Top:   5,
				Style: fontstyle.Bold,
				Size:  20,
				Align: align.Center,
			}),
		),
	)

	// Details Row
	m.AddRow(10,
		col.New(6).Add(
			text.New(fmt.Sprintf("Order ID: %d", order.ID), props.Text{Size: 12}),
		),
		col.New(6).Add(
			text.New(fmt.Sprintf("Date: %s", order.CreatedAt.Format("2006-01-02")), props.Text{
				Size:  12,
				Align: align.Right,
			}),
		),
	)

	// Price Row
	m.AddRow(15,
		col.New(12).Add(
			text.New(fmt.Sprintf("Total Estimated Cost: $%.2f", order.TotalAmountUSD), props.Text{
				Top:   5,
				Style: fontstyle.Bold,
				Size:  16,
			}),
		),
	)

	// Footer
	m.AddRow(10,
		col.New(12).Add(
			text.New("Powered by Sheet Metal Automation (Faz 6.4)", props.Text{
				Top:   5,
				Size:  10,
				Align: align.Center,
			}),
		),
	)

	// Generate PDF
	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}
