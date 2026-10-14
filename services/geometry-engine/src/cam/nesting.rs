use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct NestingRequest {
    pub sheet_width: f64,
    pub sheet_height: f64,
    pub parts: Vec<PartToNest>,
    pub allow_rotation: bool,
    pub heuristic_benchmark_ms: u64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct PartToNest {
    pub id: String,
    pub width: f64,
    pub height: f64,
    pub quantity: u32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PlacedPart {
    pub part_id: String,
    pub x: f64,
    pub y: f64,
    pub rotated: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct NestingResponse {
    pub utilized_area: f64,
    pub total_sheet_area: f64,
    pub scrap_percentage: f64,
    pub placed_parts: Vec<PlacedPart>,
    pub unplaced_count: u32,
    pub estimated_time_saved_ms: u64, 
}

#[instrument(skip(req))]
pub fn calculate_nesting(req: &NestingRequest) -> Result<NestingResponse, EngineError> {
    info!("Calculating nesting for {} distinct parts", req.parts.len());

    if req.sheet_width <= 0.0 || req.sheet_height <= 0.0 {
        return Err(EngineError::InvalidInput("Sheet dimensions must be > 0".to_string()));
    }

    let sheet_area = req.sheet_width * req.sheet_height;
    
    let mut parts = Vec::new();
    for p in &req.parts {
        for _ in 0..p.quantity {
            parts.push(crate::nesting::shapes::Part {
                id: p.id.clone(),
                width: p.width,
                height: p.height,
                placed_x: 0.0,
                placed_y: 0.0,
                rotated: false,
            });
        }
    }

    let sheets = crate::nesting::optimizer::NestingOptimizer::optimize(req.sheet_width, req.sheet_height, parts);
    
    let mut utilized_area = 0.0;
    let mut placed_parts = Vec::new();
    let mut unplaced_count = 0; 

    if !sheets.is_empty() {
        let first_sheet = &sheets[0];
        utilized_area = first_sheet.used_area();
        
        for p in &first_sheet.placed_parts {
            placed_parts.push(PlacedPart {
                part_id: p.id.clone(),
                x: p.placed_x,
                y: p.placed_y,
                rotated: p.rotated,
            });
        }
        
        let total_parts: u32 = req.parts.iter().map(|p| p.quantity).sum();
        unplaced_count = total_parts - first_sheet.placed_parts.len() as u32;
    }


    let scrap_percentage = ((sheet_area - utilized_area) / sheet_area) * 100.0;

    Ok(NestingResponse {
        utilized_area,
        total_sheet_area: sheet_area,
        scrap_percentage,
        placed_parts,
        unplaced_count,
        estimated_time_saved_ms: req.heuristic_benchmark_ms,
    })
}
