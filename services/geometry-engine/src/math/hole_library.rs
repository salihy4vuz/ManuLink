use serde::{Deserialize, Serialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct HoleRequest {
    pub metric_size: String, // e.g., "M3", "M4", "M5"
    pub hole_type: String, // "Tapped", "Clearance_Close", "Clearance_Standard"
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HoleResponse {
    pub diameter: f64,
}

#[instrument(skip(req))]
pub fn get_hole_diameter(req: &HoleRequest) -> HoleResponse {
    info!("Resolving hole diameter for {} ({})", req.metric_size, req.hole_type);
    // Basic ISO reference for Metric holes
    let dia = match req.metric_size.to_uppercase().as_str() {
        "M3" => {
            match req.hole_type.to_lowercase().as_str() {
                "tapped" => 2.5,
                "clearance_close" => 3.2,
                _ => 3.4, // standard
            }
        },
        "M4" => {
            match req.hole_type.to_lowercase().as_str() {
                "tapped" => 3.3,
                "clearance_close" => 4.3,
                _ => 4.5,
            }
        },
        "M5" => {
            match req.hole_type.to_lowercase().as_str() {
                "tapped" => 4.2,
                "clearance_close" => 5.3,
                _ => 5.5,
            }
        },
        _ => 0.0, // unknown
    };

    HoleResponse { diameter: dia }
}
