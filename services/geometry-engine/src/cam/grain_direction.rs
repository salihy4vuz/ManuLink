use serde::{Serialize, Deserialize};
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct GrainDirectionRequest {
    pub material: String,
    pub bend_angle_to_grain: f64, // 0 = Paralel, 90 = Dik
    pub thickness: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct GrainDirectionResult {
    pub is_safe: bool,
    pub recommended_action: Option<String>,
}

pub struct GrainDirectionAnalyzer;

impl GrainDirectionAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &GrainDirectionRequest) -> GrainDirectionResult {
        info!("Analyzing grain direction for material {}", req.material);

        let mat = req.material.to_lowercase();
        if mat.contains("stainless") || mat.contains("aluminum") {
            if req.bend_angle_to_grain < 45.0 || req.bend_angle_to_grain > 135.0 {
                let msg = "Bend is PARALLEL to grain. Material will crack. Rotate part by 90 degrees in nesting.";
                warn!("{}", msg);
                return GrainDirectionResult {
                    is_safe: false,
                    recommended_action: Some(msg.to_string()),
                };
            }
        }
        
        GrainDirectionResult {
            is_safe: true,
            recommended_action: None,
        }
    }
}
