use serde::{Serialize, Deserialize};
use tracing::{info, instrument};
use anyhow::Result;

#[derive(Debug, Serialize, Deserialize)]
pub struct DeburringRequest {
    pub part_area_mm2: f64,
    pub thickness_mm: f64,
    pub cut_process: String, // e.g. "laser", "plasma", "punch"
    pub plasma_multiplier: f64,
    pub punch_multiplier: f64,
    pub laser_multiplier: f64,
    pub base_time_per_10000mm2_sec: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct DeburringResult {
    pub is_deburring_required: bool,
    pub estimated_time_sec: f64,
}

pub struct DeburringAnalyzer;

impl DeburringAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &DeburringRequest) -> Result<DeburringResult> {
        info!("Analyzing deburring for process {}", req.cut_process);

        let is_required = true; // B2B parts generally require it

        let process_multiplier = match req.cut_process.to_lowercase().as_str() {
            "plasma" => req.plasma_multiplier,
            "punch" => req.punch_multiplier,
            _ => req.laser_multiplier,
        };

        let base_time_sec = (req.part_area_mm2 / 10000.0) * req.base_time_per_10000mm2_sec;
        let total_time = base_time_sec * process_multiplier;

        Ok(DeburringResult {
            is_deburring_required: is_required,
            estimated_time_sec: total_time,
        })
    }
}
