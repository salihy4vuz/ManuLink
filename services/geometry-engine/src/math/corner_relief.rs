use serde::{Deserialize, Serialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct CornerReliefRequest {
    pub relief_type: String, 
    pub thickness: f64,
    pub width_multiplier: f64,
    pub depth_multiplier: f64,
    pub radius_multiplier: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CornerReliefResponse {
    pub width: f64,
    pub depth: f64,
    pub radius: Option<f64>,
}

#[instrument(skip(req))]
pub fn calculate_corner_relief(req: &CornerReliefRequest) -> CornerReliefResponse {
    info!("Calculating corner relief for type: {}", req.relief_type);
    
    match req.relief_type.to_lowercase().as_str() {
        "circular" => CornerReliefResponse {
            width: req.thickness * req.width_multiplier,
            depth: req.thickness * req.depth_multiplier,
            radius: Some(req.thickness * req.radius_multiplier),
        },
        "rectangular" => CornerReliefResponse {
            width: req.thickness * req.width_multiplier,
            depth: req.thickness * req.depth_multiplier,
            radius: None,
        },
        "tear" | _ => CornerReliefResponse {
            width: 0.0,
            depth: 0.0,
            radius: None,
        }
    }
}
