package controllers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/salihy4vuz/sheet-metal-automation/services/main-api/rabbitmq"
)

func createGLB(jsonBytes []byte) []byte {
	padding := (4 - (len(jsonBytes) % 4)) % 4
	for i := 0; i < padding; i++ {
		jsonBytes = append(jsonBytes, ' ')
	}
	
	jsonChunkLen := uint32(len(jsonBytes))
	totalLen := 12 + 8 + jsonChunkLen
	
	buf := new(bytes.Buffer)
	buf.WriteString("glTF")
	binary.Write(buf, binary.LittleEndian, uint32(2))
	binary.Write(buf, binary.LittleEndian, totalLen)
	
	binary.Write(buf, binary.LittleEndian, jsonChunkLen)
	buf.WriteString("JSON")
	buf.Write(jsonBytes)
	
	return buf.Bytes()
}

// ViewerGLTF calls Rust engine to generate a 3D GLTF representation of the part.
// GET /api/v1/cad/:id/viewer
func ViewerGLTF(c *fiber.Ctx) error {
	fileID := c.Params("id")
	log.Printf("[Export] Requesting GLTF for %s", fileID)

	payload := map[string]interface{}{
		"part_id": fileID,
		"vertices": [][]float32{
			{0.0, 0.0, 0.0}, {100.0, 0.0, 0.0}, {100.0, 100.0, 0.0}, {0.0, 100.0, 0.0},
		},
		"indices": []uint16{0, 1, 2, 0, 2, 3},
	}

	res, err := rabbitmq.CallRustEngine(c.Context(), "export_gltf", payload)
	if err != nil {
		log.Printf("[Export] Error calling Rust Engine for GLTF: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate 3D viewer file",
		})
	}

	glbBytes := createGLB([]byte(res))

	c.Set("Content-Type", "model/gltf-binary")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.glb"`, fileID))
	c.Set("Cache-Control", "public, max-age=3600")
	return c.Send(glbBytes)
}

// DownloadDXF calls Rust engine to generate a 2D DXF for laser/bending.
// GET /api/v1/cad/:id/download-dxf
func DownloadDXF(c *fiber.Ctx) error {
	fileID := c.Params("id")
	log.Printf("[Export] Requesting DXF for %s", fileID)

	payload := map[string]interface{}{
		"layer_name": "CUT_OUTER",
		"lines": []map[string]float64{
			{"x1": 0.0, "y1": 0.0, "x2": 100.0, "y2": 0.0},
			{"x1": 100.0, "y1": 0.0, "x2": 100.0, "y2": 100.0},
			{"x1": 100.0, "y1": 100.0, "x2": 0.0, "y2": 100.0},
			{"x1": 0.0, "y1": 100.0, "x2": 0.0, "y2": 0.0},
		},
		"circles": []map[string]float64{},
	}

	res, err := rabbitmq.CallRustEngine(c.Context(), "export_dxf", payload)
	if err != nil {
		log.Printf("[Export] Error calling Rust Engine for DXF: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate DXF",
		})
	}

	c.Set("Content-Type", "application/dxf")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.dxf"`, fileID))
	return c.SendString(res)
}
