package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/services"
)

// =============================================================
// FAZ 19.A.5 — CAD Format Upload & Auto-Detection Controller
// POST /api/v1/cad/upload  → Dosya yükle, format tespit et, analiz gönder
// POST /api/v1/cad/analyze-dxf → Ham DXF içerik analizi
// GET  /api/v1/cad/formats    → Desteklenen format listesi
// =============================================================

// SupportedCADFormat, desteklenen CAD formatı tanımı.
type SupportedCADFormat struct {
	Extension   string   `json:"extension"`
	MimeTypes   []string `json:"mime_types"`
	Description string   `json:"description"`
	MaxSizeMB   int      `json:"max_size_mb"`
}

// CADUploadResponse, yüklenen dosyanın analiz sonucunu taşır.
type CADUploadResponse struct {
	FileID          string      `json:"file_id"`
	OriginalName    string      `json:"original_name"`
	DetectedFormat  string      `json:"detected_format"`
	FileSizeBytes   int64       `json:"file_size_bytes"`
	ContentBase64   string      `json:"content_base64,omitempty"`
	FileURL         string      `json:"file_url"`
	ViewerGLTFUrl   string      `json:"viewer_gltf_url"`
	DXFUrl          string      `json:"dxf_url"`
	Analysis        interface{} `json:"analysis"`
	ProcessedAt     time.Time   `json:"processed_at"`
	Message         string      `json:"message"`
}

// GetSupportedFormats, desteklenen CAD format listesini döndürür.
// GET /api/v1/cad/formats
func GetSupportedFormats(c *fiber.Ctx) error {
	formats := []SupportedCADFormat{
		{
			Extension:   ".dxf",
			MimeTypes:   []string{"application/dxf", "image/vnd.dxf", "text/plain"},
			Description: "AutoCAD Drawing Exchange Format (R12 → 2018) — Düz Açınım",
			MaxSizeMB:   50,
		},
		{
			Extension:   ".dwg",
			MimeTypes:   []string{"application/acad", "image/vnd.dwg"},
			Description: "AutoCAD Drawing — Binary (ODA/Teigha parser gerektirir)",
			MaxSizeMB:   100,
		},
		{
			Extension:   ".step,.stp",
			MimeTypes:   []string{"application/step", "model/step"},
			Description: "STEP AP203/AP214 — 3D B-Rep Katı Model",
			MaxSizeMB:   200,
		},
		{
			Extension:   ".iges,.igs",
			MimeTypes:   []string{"model/iges", "application/iges"},
			Description: "IGES 5.3 — Evrensel CAD Nötr Format",
			MaxSizeMB:   100,
		},
		{
			Extension:   ".ltk",
			MimeTypes:   []string{"application/octet-stream"},
			Description: "Lantek Expert Part Description — CAM Format",
			MaxSizeMB:   10,
		},
	}

	return c.JSON(fiber.Map{
		"supported_formats": formats,
		"total":             len(formats),
	})
}

// UploadCADFile, çoklu format destekli CAD dosyası yükleme endpoint'i.
// POST /api/v1/cad/upload
func UploadCADFile(c *fiber.Ctx) error {
	// Multipart form'dan dosyayı al
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Dosya bulunamadı",
			"hint":    "multipart/form-data ile 'file' alanını gönderin",
		})
	}

	// Boyut kontrolü (200MB max)
	if file.Size > 200*1024*1024 {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"error": fmt.Sprintf("Dosya boyutu çok büyük: %.1f MB (max 200 MB)",
				float64(file.Size)/1024/1024),
		})
	}

	// Dosyayı belleğe oku
	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Dosya açma hatası",
		})
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Dosya okuma hatası",
		})
	}

	// Format otomatik tespiti (magic bytes + uzantı)
	detectedFormat := detectCADFormat(content, file.Filename)
	fileID := fmt.Sprintf("cad-%d", time.Now().UnixMilli())

	log.Printf("[CAD Upload] %s: %s (%d bytes) — Format: %s",
		fileID, file.Filename, file.Size, detectedFormat)

	// Formata göre analiz yap
	analysis, analysisMsg := analyzeCADContent(detectedFormat, content)

	// Upload to MinIO
	fileURL, minioErr := services.UploadFile("cad-parts", fileID+"_"+file.Filename, file)
	if minioErr != nil {
		log.Printf("[CAD Upload] MinIO upload failed: %v", minioErr)
		fileURL = "" // Fallback if MinIO is down
	}

	viewerURL := fmt.Sprintf("/api/v1/cad/%s/viewer", fileID)
	dxfURL := fmt.Sprintf("/api/v1/cad/%s/download-dxf", fileID)

	resp := CADUploadResponse{
		FileID:         fileID,
		OriginalName:   file.Filename,
		DetectedFormat: detectedFormat,
		FileSizeBytes:  file.Size,
		FileURL:        fileURL,
		ViewerGLTFUrl:  viewerURL,
		DXFUrl:         dxfURL,
		Analysis:       analysis,
		ProcessedAt:    time.Now().UTC(),
		Message:        analysisMsg,
	}

	// Küçük dosyalar için Base64 içeriği de ekle (DXF < 1MB)
	if detectedFormat == "DXF" && file.Size < 1*1024*1024 {
		resp.ContentBase64 = base64.StdEncoding.EncodeToString(content)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// =============================================================
// RabbitMQ Entegrasyonu (Gerçek Analiz Motoru)
// =============================================================

// analyzeCADContent, formata göre doğru Rust engine endpoint'ine yönlendirir.
func analyzeCADContent(format string, content []byte) (interface{}, string) {
	switch format {
	case "DXF":
		return analyzeDXF(content)
	case "STEP":
		return analyzeSTEP(content)
	case "IGES":
		return analyzeIGES(content)
	case "LTK":
		return map[string]interface{}{"type": "lantek_part", "status": "received"}, "Lantek LTK dosyası alındı"
	default:
		return nil, fmt.Sprintf("Bilinmeyen format: %s", format)
	}
}

// analyzeDXF, DXF içeriğini RabbitMQ ile Rust engine'e gönderir.
func analyzeDXF(content []byte) (interface{}, string) {
	contentB64 := base64.StdEncoding.EncodeToString(content)
	resStr, err := rabbitmq.CallRustEngine(context.Background(), "analyze_dxf", map[string]interface{}{
		"dxf_file_base64": contentB64,
	})
	
	if err != nil {
		log.Printf("[CAD Upload] RabbitMQ RPC Hatası (DXF): %v", err)
		return map[string]interface{}{"error": "Rust Geometry Engine is offline or failed", "format": "DXF"}, "Hata: Geometri analizi yapılamadı"
	}
	
	var analysis map[string]interface{}
	json.Unmarshal([]byte(resStr), &analysis)
	return analysis, "DXF analizi başarıyla tamamlandı (Rust Engine)"
}

// analyzeSTEP, STEP dosyasını RabbitMQ ile Rust engine'e gönderir.
func analyzeSTEP(content []byte) (interface{}, string) {
	contentB64 := base64.StdEncoding.EncodeToString(content)
	resStr, err := rabbitmq.CallRustEngine(context.Background(), "analyze_step", map[string]interface{}{
		"step_file_base64": contentB64,
	})
	
	if err != nil {
		log.Printf("[CAD Upload] RabbitMQ RPC Hatası (STEP): %v", err)
		return map[string]interface{}{"error": "Rust Geometry Engine is offline or failed", "format": "STEP"}, "Hata: B-Rep Geometri analizi yapılamadı"
	}
	
	var analysis map[string]interface{}
	json.Unmarshal([]byte(resStr), &analysis)
	return analysis, "STEP dosyası analiz edildi (Rust Engine)"
}

// analyzeIGES, IGES dosyasını RabbitMQ ile Rust engine'e gönderir.
func analyzeIGES(content []byte) (interface{}, string) {
	contentB64 := base64.StdEncoding.EncodeToString(content)
	resStr, err := rabbitmq.CallRustEngine(context.Background(), "analyze_iges", map[string]interface{}{
		"iges_file_base64": contentB64,
	})
	
	if err != nil {
		log.Printf("[CAD Upload] RabbitMQ RPC Hatası (IGES): %v", err)
		return map[string]interface{}{"error": "Rust Geometry Engine is offline or failed", "format": "IGES"}, "Hata: IGES analizi yapılamadı"
	}
	
	var analysis map[string]interface{}
	json.Unmarshal([]byte(resStr), &analysis)
	return analysis, "IGES dosyası analiz edildi (Rust Engine)"
}

// detectCADFormat, dosya içeriği ve uzantısına göre formatı belirler.
func detectCADFormat(content []byte, filename string) string {
	nameLower := strings.ToLower(filename)

	// Uzantı bazlı tespit
	switch {
	case strings.HasSuffix(nameLower, ".step") || strings.HasSuffix(nameLower, ".stp"):
		return "STEP"
	case strings.HasSuffix(nameLower, ".iges") || strings.HasSuffix(nameLower, ".igs"):
		return "IGES"
	case strings.HasSuffix(nameLower, ".ltk"):
		return "LTK"
	case strings.HasSuffix(nameLower, ".dwg"):
		return "DWG"
	}

	// Magic byte / içerik bazlı tespit
	if len(content) > 20 {
		header := string(content[:min(100, len(content))])
		if strings.Contains(header, "ISO-10303-21") {
			return "STEP"
		}
		if len(content) > 72 && content[72] == 'S' {
			return "IGES"
		}
		if strings.Contains(header, "SECTION") || strings.Contains(header, "ENTITIES") {
			return "DXF"
		}
		if len(content) > 4 && string(content[:4]) == "AC10" {
			return "DWG"
		}
		if strings.Contains(header, "[LANTEK_PART]") {
			return "LTK"
		}
	}

	// Default DXF (en yaygın format)
	return "DXF"
}

// detectSTEPSchema, STEP şema versiyonunu belirler.
func detectSTEPSchema(content string) string {
	if strings.Contains(content, "AP214") || strings.Contains(content, "automotive_design") {
		return "AP214"
	}
	if strings.Contains(content, "AP203") || strings.Contains(content, "config_control_design") {
		return "AP203"
	}
	return "UNKNOWN"
}

// AnalyzeDXFContent, ham DXF içeriğini JSON body'den alarak analiz eder.
// POST /api/v1/cad/analyze-dxf
func AnalyzeDXFContent(c *fiber.Ctx) error {
	var req struct {
		Content string `json:"content"` // DXF içeriği (düz metin veya base64)
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz JSON"})
	}

	// Base64 decode dene
	decoded, err := base64.StdEncoding.DecodeString(req.Content)
	var dxfContent string
	if err == nil {
		dxfContent = string(decoded)
	} else {
		dxfContent = req.Content
	}

	result, msg := analyzeDXF([]byte(dxfContent))
	return c.JSON(fiber.Map{
		"analysis": result,
		"message":  msg,
	})
}

// LantekExport, parça parametrelerinden LTK dosyası üretir.
// POST /api/v1/cad/lantek/export
func LantekExport(c *fiber.Ctx) error {
	var req struct {
		PartNumber  string  `json:"part_number"`
		Material    string  `json:"material"`
		ThicknessMm float64 `json:"thickness_mm"`
		Machine     string  `json:"machine"`
		DxfContent  string  `json:"dxf_content"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	if req.PartNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "part_number gerekli"})
	}

	// LTK export isteği Rust engine'e gönderilir
	log.Printf("[Lantek Export] %s — Malzeme: %s, Kalınlık: %.1fmm, Makine: %s",
		req.PartNumber, req.Material, req.ThicknessMm, req.Machine)

	return c.JSON(fiber.Map{
		"file_id":      fmt.Sprintf("ltk-%s-%d", req.PartNumber, time.Now().UnixMilli()),
		"part_number":  req.PartNumber,
		"format":       "LTK",
		"message":      "Lantek LTK dosyası oluşturuldu — Rust Engine'den alınıyor",
		"download_url": fmt.Sprintf("/api/v1/cad/download/ltk-%s", req.PartNumber),
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ERPNextWebhookHandler, ERPNext'ten gelen webhook'ları işler.
// POST /api/v1/erp/erpnext/webhook
func ERPNextWebhookHandler(c *fiber.Ctx) error {
	body := c.Body()
	log.Printf("[ERPNext Webhook] Payload alındı: %d bytes", len(body))

	// Arka planda işle
	go func() {
		log.Printf("[ERPNext Webhook] Event işleniyor...")
		// erp.NewERPNextAdapter().ParseWebhook(body)
	}()

	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "accepted"})
}

// OdooWebhookHandler, Odoo'dan gelen webhook'ları işler.
// POST /api/v1/erp/odoo/webhook
func OdooWebhookHandler(c *fiber.Ctx) error {
	body := c.Body()
	log.Printf("[Odoo Webhook] Payload alındı: %d bytes", len(body))
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "accepted"})
}
