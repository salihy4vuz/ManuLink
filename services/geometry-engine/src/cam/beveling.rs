use serde::{Serialize, Deserialize};
use tracing::{info, instrument};
use anyhow::Result;

#[derive(Debug, Serialize, Deserialize)]
pub enum BevelType {
    VBevel,
    YBevel,
    XBevel,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct BevelRequest {
    pub thickness_mm: f64,
    pub edge_length_mm: f64,
    pub bevel_type: BevelType,
    pub required_thickness_threshold: f64,
    pub material_removal_rate_mm3_per_sec: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct BevelResult {
    pub estimated_machining_time_sec: f64,
    pub volume_removed_mm3: f64,
    pub is_required: bool,
}

pub struct BevelAnalyzer;

impl BevelAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &BevelRequest) -> Result<BevelResult> {
        info!("Analyzing beveling requirements for edge length {}", req.edge_length_mm);

        let is_required = req.thickness_mm >= req.required_thickness_threshold;

        let cross_section_area = match req.bevel_type {
            BevelType::VBevel => (req.thickness_mm * req.thickness_mm) / 2.0,
            BevelType::YBevel => ((req.thickness_mm - 2.0).max(0.0) * req.thickness_mm) / 2.0,
            BevelType::XBevel => (req.thickness_mm * req.thickness_mm) / 4.0,
        };

        let volume_removed = cross_section_area * req.edge_length_mm;
        let time_sec = volume_removed / req.material_removal_rate_mm3_per_sec;

        Ok(BevelResult {
            estimated_machining_time_sec: time_sec,
            volume_removed_mm3: volume_removed,
            is_required,
        })
    }
}
