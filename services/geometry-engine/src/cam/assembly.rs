use serde::{Serialize, Deserialize};
use tracing::{info, instrument};
use anyhow::Result;

#[derive(Debug, Serialize, Deserialize)]
pub struct AssemblyRequest {
    pub parts_count: u32,
    pub total_weld_length_mm: f64,
    pub requires_custom_jig: bool,
    // Dynamic config params
    pub base_jig_setup_sec: f64,
    pub per_part_jig_setup_sec: f64,
    pub per_part_standard_setup_sec: f64,
    pub welding_speed_mm_per_sec: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AssemblyResult {
    pub estimated_jig_setup_time_sec: f64,
    pub estimated_welding_time_sec: f64,
    pub total_assembly_time_sec: f64,
}

pub struct AssemblyAnalyzer;

impl AssemblyAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &AssemblyRequest) -> Result<AssemblyResult> {
        info!("Analyzing assembly time for {} parts", req.parts_count);

        let jig_setup = if req.requires_custom_jig {
            req.base_jig_setup_sec + (req.per_part_jig_setup_sec * req.parts_count as f64)
        } else {
            req.per_part_standard_setup_sec * req.parts_count as f64
        };

        let welding_time = req.total_weld_length_mm / req.welding_speed_mm_per_sec;

        Ok(AssemblyResult {
            estimated_jig_setup_time_sec: jig_setup,
            estimated_welding_time_sec: welding_time,
            total_assembly_time_sec: jig_setup + welding_time,
        })
    }
}
