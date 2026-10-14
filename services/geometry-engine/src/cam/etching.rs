use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct EtchingRequest {
    pub text_length_mm: f64,
    pub line_count: u32,
    pub speed_mm_sec: f64, 
    pub jump_distance_multiplier: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct EtchingResult {
    pub estimated_time_sec: f64,
}

pub struct EtchingAnalyzer;

impl EtchingAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &EtchingRequest) -> EtchingResult {
        info!("Analyzing etching time for {} lines", req.line_count);
        
        let total_distance = req.text_length_mm * req.line_count as f64 * req.jump_distance_multiplier; 

        EtchingResult {
            estimated_time_sec: total_distance / req.speed_mm_sec,
        }
    }
}
