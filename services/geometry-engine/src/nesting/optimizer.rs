use super::maxrects::{MaxRectsPacker, FreeRectChoiceHeuristic};
use super::shapes::{Part, Sheet};
use tracing::{warn, instrument};

pub struct NestingOptimizer;

impl NestingOptimizer {
    /// Attempts to fit all parts on the given sheets.
    #[instrument(skip(parts))]
    pub fn optimize(sheet_width: f64, sheet_height: f64, mut parts: Vec<Part>) -> Vec<Sheet> {
        // Sort parts by area (largest first) as a simple heuristic
        parts.sort_by(|a, b| {
            let area_b = b.width * b.height;
            let area_a = a.width * a.height;
            area_b.partial_cmp(&area_a).unwrap()
        });

        let mut sheets: Vec<Sheet> = Vec::new();
        let packer = MaxRectsPacker::new();
        
        let mut current_sheet = Sheet::new(sheet_width, sheet_height);

        for part in parts {
            let mut placed = packer.insert_part(&mut current_sheet, part.clone(), FreeRectChoiceHeuristic::BestShortSideFit);
            
            if !placed {
                // If it doesn't fit on current sheet, store current sheet and create a new one
                sheets.push(current_sheet);
                current_sheet = Sheet::new(sheet_width, sheet_height);
                placed = packer.insert_part(&mut current_sheet, part.clone(), FreeRectChoiceHeuristic::BestShortSideFit);
                
                if !placed {
                    warn!("Warning: Part {} is larger than the sheet size itself!", part.id);
                }
            }
        }
        
        // Add the last sheet
        if current_sheet.used_area() > 0.0 {
            sheets.push(current_sheet);
        }

        sheets
    }
}
