package services

import (
	"math"
)

// Standard Sheet Sizes (in mm)
const (
	SheetLength2000 = 2000.0
	SheetWidth1000  = 1000.0

	SheetLength2500 = 2500.0
	SheetWidth1250  = 1250.0

	SheetLength3000 = 3000.0
	SheetWidth1500  = 1500.0
)

type SheetRequirement struct {
	SheetLength float64
	SheetWidth  float64
	SheetsNeeded int
	ScrapRate    float64
	TotalCost    float64
}

// CalculateSheetUtilization determines the most efficient standard sheet to use
// and how many of them are needed to fulfill the quantity of parts, including 10mm margins.
func CalculateSheetUtilization(partX, partY float64, quantity int, materialKgPrice float64, thickness float64, density float64) SheetRequirement {
	// Add 10mm margin between parts for laser cutting
	paddedX := partX + 10.0
	paddedY := partY + 10.0

	// We will try 3 standard sheet sizes and find the one with the lowest scrap rate
	sheets := []struct {
		L float64
		W float64
	}{
		{SheetLength2000, SheetWidth1000},
		{SheetLength2500, SheetWidth1250},
		{SheetLength3000, SheetWidth1500},
	}

	var bestReq SheetRequirement
	bestReq.ScrapRate = 100.0 // Max initial scrap

	for _, sheet := range sheets {
		// Calculate how many parts fit in one sheet (simple grid nesting)
		partsPerRow := math.Floor(sheet.L / paddedX)
		partsPerCol := math.Floor(sheet.W / paddedY)
		
		// Also try rotated
		partsPerRowRotated := math.Floor(sheet.L / paddedY)
		partsPerColRotated := math.Floor(sheet.W / paddedX)

		partsPerSheet1 := partsPerRow * partsPerCol
		partsPerSheet2 := partsPerRowRotated * partsPerColRotated

		partsPerSheet := math.Max(partsPerSheet1, partsPerSheet2)

		if partsPerSheet > 0 {
			sheetsNeeded := int(math.Ceil(float64(quantity) / partsPerSheet))
			
			// Calculate Scrap Rate
			totalPartArea := (partX * partY) * float64(quantity)
			totalSheetArea := (sheet.L * sheet.W) * float64(sheetsNeeded)
			scrapRate := ((totalSheetArea - totalPartArea) / totalSheetArea) * 100.0

			// Calculate cost of the entire plates bought from supplier
			// Volume in mm^3 -> cm^3 (divide by 1000)
			volCubicCm := (sheet.L * sheet.W * thickness) / 1000.0
			massKg := volCubicCm * density
			totalPlateCost := massKg * materialKgPrice * float64(sheetsNeeded)

			// We choose the sheet that gives the lowest scrap rate
			if scrapRate < bestReq.ScrapRate {
				bestReq = SheetRequirement{
					SheetLength:  sheet.L,
					SheetWidth:   sheet.W,
					SheetsNeeded: sheetsNeeded,
					ScrapRate:    scrapRate,
					TotalCost:    totalPlateCost,
				}
			}
		}
	}

	return bestReq
}
