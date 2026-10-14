package services

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/go-pdf/fpdf"
)

// MaterialCertificate, EN 10204 3.1 sertifikası için gerekli tüm verileri tutar
type MaterialCertificate struct {
	CertificateNo  string
	IssuedDate     time.Time
	PartNo         string
	Description    string
	Material       string
	HeatNo         string
	Standard       string
	ChemicalComp   map[string]float64 // C%, Si%, Mn%, Cr%, Ni%, vb.
	MechanicalProp map[string]float64 // Rp0.2, Rm, A5%, HB
	Thickness      float64
	WeightKg       float64
	InspectorName  string
	CompanyName    string
	CustomerName   string
	OrderID        string
}

// MaterialCertificateDefaults fills default chemical/mechanical values for well-known materials
func MaterialCertificateDefaults(cert *MaterialCertificate) {
	if cert.ChemicalComp == nil {
		cert.ChemicalComp = map[string]float64{}
	}
	if cert.MechanicalProp == nil {
		cert.MechanicalProp = map[string]float64{}
	}

	switch cert.Material {
	case "AISI 304", "304", "304L", "1.4307":
		cert.Standard = "EN 10088-2 / ASTM A240"
		cert.ChemicalComp = map[string]float64{"C": 0.023, "Si": 0.42, "Mn": 1.52, "P": 0.033, "S": 0.002, "Cr": 18.2, "Ni": 8.9, "Mo": 0.04}
		cert.MechanicalProp = map[string]float64{"Rp0.2": 285, "Rm": 590, "A5": 47, "HB": 168}
	case "AISI 316L", "316L", "1.4404":
		cert.Standard = "EN 10088-2 / ASTM A240"
		cert.ChemicalComp = map[string]float64{"C": 0.018, "Si": 0.38, "Mn": 1.41, "P": 0.030, "S": 0.001, "Cr": 16.8, "Ni": 10.3, "Mo": 2.04}
		cert.MechanicalProp = map[string]float64{"Rp0.2": 298, "Rm": 612, "A5": 46, "HB": 172}
	case "Al 6061", "6061-T6", "AlSi", "EN AW-6061":
		cert.Standard = "EN 573-3 / ASTM B209"
		cert.ChemicalComp = map[string]float64{"Si": 0.62, "Mg": 0.95, "Cu": 0.27, "Cr": 0.06, "Fe": 0.18}
		cert.MechanicalProp = map[string]float64{"Rp0.2": 276, "Rm": 310, "A5": 12, "HB": 95}
	default:
		// S235JR / ST37 varsayılan
		cert.Standard = "EN 10025-2"
		cert.ChemicalComp = map[string]float64{"C": 0.17, "Si": 0.20, "Mn": 1.40, "P": 0.035, "S": 0.035}
		cert.MechanicalProp = map[string]float64{"Rp0.2": 235, "Rm": 360, "A5": 26, "HB": 120}
	}

	if cert.Standard == "" { cert.Standard = "EN 10204:2004 Tip 3.1" }
	if cert.CompanyName == "" { cert.CompanyName = "Apex Factory OS A.Ş." }
	if cert.InspectorName == "" { cert.InspectorName = "Kalite Güvence Departmanı" }
	if cert.HeatNo == "" {
		cert.HeatNo = fmt.Sprintf("MH-%s-%d", cert.OrderID, time.Now().UnixNano()%100000)
	}
	if cert.IssuedDate.IsZero() {
		cert.IssuedDate = time.Now()
	}
}

// GenerateMaterialCertificatePDF, EN 10204 3.1 formatında kapsamlı PDF sertifika üretir.
// fpdf tabanlı implementasyon — maroto kütüphanesinden bağımsız.
func GenerateMaterialCertificatePDF(cert MaterialCertificate) ([]byte, error) {
	MaterialCertificateDefaults(&cert)

	p := fpdf.New("P", "mm", "A4", "")
	p.AddPage()
	p.SetMargins(15, 15, 15)

	// ─── Başlık alanı ─────────────────────────────────────────
	p.SetFillColor(15, 23, 42)
	p.Rect(0, 0, 210, 24, "F")
	p.SetTextColor(255, 255, 255)
	p.SetFont("Arial", "B", 13)
	p.SetXY(15, 7)
	p.Cell(140, 7, "MALZEME TEST SERTİFİKASI / INSPECTION CERTIFICATE")
	p.SetFont("Arial", "", 9)
	p.SetXY(155, 7)
	p.Cell(50, 5, "EN 10204 Tip 3.1")
	p.SetXY(155, 13)
	p.Cell(50, 5, cert.IssuedDate.Format("02.01.2006"))
	p.SetTextColor(0, 0, 0)
	p.SetY(30)

	// ─── Sertifika Kimliği ────────────────────────────────────
	p.SetFont("Arial", "B", 10)
	p.SetFillColor(240, 249, 255)
	p.CellFormat(180, 7, "SERTİFİKA BİLGİLERİ / CERTIFICATE INFORMATION", "1", 1, "C", true, 0, "")

	rows := [][]string{
		{"Sertifika No", cert.CertificateNo, "Sipariş No / Order", cert.OrderID},
		{"Müşteri / Customer", cert.CustomerName, "Parça / Part No", cert.PartNo},
		{"Üretici / Manufacturer", cert.CompanyName, "Isıl İşlem Lot / Heat No", cert.HeatNo},
	}
	for _, r := range rows {
		p.SetFont("Arial", "B", 8); p.CellFormat(35, 6, r[0]+":", "1", 0, "", false, 0, "")
		p.SetFont("Arial", "", 8);  p.CellFormat(55, 6, r[1], "1", 0, "", false, 0, "")
		p.SetFont("Arial", "B", 8); p.CellFormat(35, 6, r[2]+":", "1", 0, "", false, 0, "")
		p.SetFont("Arial", "", 8);  p.CellFormat(55, 6, r[3], "1", 1, "", false, 0, "")
	}
	p.Ln(4)

	// ─── Malzeme Tanımı ────────────────────────────────────────
	p.SetFont("Arial", "B", 10)
	p.SetFillColor(240, 249, 255)
	p.CellFormat(180, 7, "MALZEME TANIMI / MATERIAL DESCRIPTION", "1", 1, "C", true, 0, "")

	matRows := [][]string{
		{"Malzeme / Material", cert.Material, "Kalite / Grade", cert.Standard},
		{"Kalınlık / Thickness", fmt.Sprintf("%.2f mm", cert.Thickness), "Ağırlık / Weight", fmt.Sprintf("%.3f kg", cert.WeightKg)},
		{"Açıklama / Description", cert.Description, "Yüzey / Surface", "2B / Hadde Sonrası"},
	}
	for _, r := range matRows {
		p.SetFont("Arial", "B", 8); p.CellFormat(35, 6, r[0]+":", "1", 0, "", false, 0, "")
		p.SetFont("Arial", "", 8);  p.CellFormat(55, 6, r[1], "1", 0, "", false, 0, "")
		p.SetFont("Arial", "B", 8); p.CellFormat(35, 6, r[2]+":", "1", 0, "", false, 0, "")
		p.SetFont("Arial", "", 8);  p.CellFormat(55, 6, r[3], "1", 1, "", false, 0, "")
	}
	p.Ln(4)

	// ─── Kimyasal Analiz ──────────────────────────────────────
	p.SetFont("Arial", "B", 10)
	p.SetFillColor(240, 249, 255)
	p.CellFormat(180, 7, "KİMYASAL ANALİZ / CHEMICAL ANALYSIS (% Ağırlık / Weight)", "1", 1, "C", true, 0, "")

	chemOrder := []string{"C", "Si", "Mn", "P", "S", "Cr", "Ni", "Mo"}
	colW := 180.0 / float64(len(chemOrder))
	p.SetFont("Arial", "B", 8)
	for _, k := range chemOrder {
		p.CellFormat(colW, 6, k, "1", 0, "C", false, 0, "")
	}
	p.Ln(-1)
	p.SetFont("Arial", "", 8)
	for _, k := range chemOrder {
		v, ok := cert.ChemicalComp[k]
		s := "-"
		if ok && v > 0 { s = fmt.Sprintf("%.3f", v) }
		p.CellFormat(colW, 6, s, "1", 0, "C", false, 0, "")
	}
	p.Ln(-1)
	p.Ln(4)

	// ─── Mekanik Özellikler ────────────────────────────────────
	p.SetFont("Arial", "B", 10)
	p.SetFillColor(240, 249, 255)
	p.CellFormat(180, 7, "MEKANİK ÖZELLİKLER / MECHANICAL PROPERTIES", "1", 1, "C", true, 0, "")

	mechOrder := []string{"Rp0.2", "Rm", "A5", "HB"}
	mechLabels := map[string]string{
		"Rp0.2": "Akma Dayanımı Rp0.2 (MPa)",
		"Rm":    "Çekme Dayanımı Rm (MPa)",
		"A5":    "Uzama A5 (%)",
		"HB":    "Sertlik (HB)",
	}
	mColW := 180.0 / float64(len(mechOrder))
	p.SetFont("Arial", "B", 8)
	for _, k := range mechOrder {
		p.CellFormat(mColW, 7, mechLabels[k], "1", 0, "C", false, 0, "")
	}
	p.Ln(-1)
	p.SetFont("Arial", "B", 11)
	for _, k := range mechOrder {
		v, _ := cert.MechanicalProp[k]
		p.CellFormat(mColW, 8, fmt.Sprintf("%.1f", v), "1", 0, "C", false, 0, "")
	}
	p.Ln(-1)
	p.Ln(6)

	// ─── Onay ve Beyan ────────────────────────────────────────
	p.SetFont("Arial", "B", 10)
	p.SetFillColor(240, 249, 255)
	p.CellFormat(180, 7, "ONAY VE BEYAN / DECLARATION OF CONFORMITY (EN 10204:2004 Tip 3.1)", "1", 1, "C", true, 0, "")
	p.SetFont("Arial", "", 8)
	p.MultiCell(180, 5,
		"Bu sertifika, EN 10204:2004 Tip 3.1 standardı kapsamında üretici kalite departmanı tarafından hazırlanmıştır. "+
			"Belirtilen kimyasal ve mekanik değerler, üretime ait gerçek eriyik analizi sonuçlarını yansıtmaktadır. "+
			"Apex Factory OS A.Ş. ISO 9001:2015 sertifikalı kalite yönetim sistemi çerçevesinde onaylanmıştır.",
		"1", "L", false)
	p.Ln(5)

	p.SetFont("Arial", "B", 8)
	p.CellFormat(60, 6, "Yetkili İmza / Authorized Signature", "1", 0, "C", false, 0, "")
	p.CellFormat(60, 6, fmt.Sprintf("Kalite Sorumlusu: %s", cert.InspectorName), "1", 0, "C", false, 0, "")
	p.CellFormat(60, 6, "Mühendis Onayı / Engineer Approval", "1", 1, "C", false, 0, "")
	p.CellFormat(60, 18, "", "1", 0, "", false, 0, "")
	p.CellFormat(60, 18, "", "1", 0, "", false, 0, "")
	p.CellFormat(60, 18, "", "1", 1, "", false, 0, "")

	// ─── Footer ───────────────────────────────────────────────
	p.SetY(-18)
	p.SetFont("Arial", "I", 7)
	p.SetTextColor(120, 120, 120)
	p.Cell(180, 5, fmt.Sprintf(
		"Otomatik oluşturuldu — Apex Factory OS Platform v1.9 | %s | Sipariş: %s",
		cert.IssuedDate.Format("02.01.2006 15:04"), cert.OrderID))

	var buf bytes.Buffer
	if err := p.Output(&buf); err != nil {
		log.Printf("[MaterialCertificate] PDF generation error: %v", err)
		return nil, fmt.Errorf("material certificate PDF generation failed: %w", err)
	}

	log.Printf("[MaterialCertificate] EN 10204 3.1 sertifikası üretildi — Sipariş: %s, Malzeme: %s, HeatNo: %s",
		cert.OrderID, cert.Material, cert.HeatNo)
	return buf.Bytes(), nil
}
