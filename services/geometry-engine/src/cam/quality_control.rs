use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct QualityControlRequest {
    pub max_dimension_mm: f64,
    pub thickness_mm: f64,
    pub tolerance_class: String, // ISO 2768-m (medium), -f (fine), -c (coarse)
    pub tight_tolerance_threshold_mm: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct QualityControlResult {
    pub required_inspection: bool,
    pub acceptable_tolerance_mm: f64,
    pub inspection_method: String,
}

pub struct QCAnalyzer;

impl QCAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &QualityControlRequest) -> QualityControlResult {
        info!("Analyzing quality control for dimension {} mm", req.max_dimension_mm);

        let mut tolerance = 0.5;

        let tc = req.tolerance_class.to_lowercase();
        if req.max_dimension_mm <= 120.0 {
            tolerance = if tc == "f" { 0.15 } else if tc == "m" { 0.3 } else { 0.8 };
        } else if req.max_dimension_mm <= 400.0 {
            tolerance = if tc == "f" { 0.2 } else if tc == "m" { 0.5 } else { 1.2 };
        } else if req.max_dimension_mm <= 1000.0 {
            tolerance = if tc == "f" { 0.3 } else if tc == "m" { 0.8 } else { 2.0 };
        }

        let is_tight = tolerance <= req.tight_tolerance_threshold_mm;
        let inspection_method = if is_tight {
            "CMM (Coordinate Measuring Machine) & Vernier Caliper".to_string()
        } else {
            "Vernier Caliper & Tape Measure".to_string()
        };

        QualityControlResult {
            required_inspection: true,
            acceptable_tolerance_mm: tolerance,
            inspection_method,
        }
    }
}
