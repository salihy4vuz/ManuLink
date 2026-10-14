use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct PaintingRequest {
    pub surface_area_mm2: f64,
    pub paint_type: String, 
    pub is_double_sided: bool,
    pub powder_consumption_g_per_m2: f64,
    pub wet_consumption_g_per_m2: f64,
    pub painting_time_sec_per_m2: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PaintingResult {
    pub total_area_m2: f64,
    pub paint_consumption_grams: f64,
    pub estimated_painting_time_sec: f64,
}

pub struct PaintingAnalyzer;

impl PaintingAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &PaintingRequest) -> PaintingResult {
        info!("Analyzing painting for surface area {} mm2", req.surface_area_mm2);

        let mut area_m2 = req.surface_area_mm2 / 1_000_000.0;
        
        if req.is_double_sided {
            area_m2 *= 2.0;
        }

        let consumption = if req.paint_type.to_lowercase() == "powder" {
            area_m2 * req.powder_consumption_g_per_m2
        } else {
            area_m2 * req.wet_consumption_g_per_m2 
        };

        let time_sec = area_m2 * req.painting_time_sec_per_m2;

        PaintingResult {
            total_area_m2: area_m2,
            paint_consumption_grams: consumption,
            estimated_painting_time_sec: time_sec,
        }
    }
}
