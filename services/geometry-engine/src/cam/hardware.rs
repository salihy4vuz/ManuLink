use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub enum HardwareType {
    Tapping, 
    PemNut,  
    ClinchStud,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HardwareRequest {
    pub hw_type: HardwareType,
    pub count: u32,
    pub thread_size: String,
    pub time_per_hole_sec: f64, // Database config provided
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HardwareResult {
    pub estimated_time_sec: f64,
}

pub struct HardwareAnalyzer;

impl HardwareAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &HardwareRequest) -> HardwareResult {
        info!("Analyzing hardware {:?} (count: {})", req.hw_type, req.count);

        HardwareResult {
            estimated_time_sec: req.time_per_hole_sec * req.count as f64,
        }
    }
}
