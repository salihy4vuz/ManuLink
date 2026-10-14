use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct CollisionRequest {
    pub flange_length: f64,
    pub bend_angle: f64,
    pub punch_height: f64,
    pub v_die_width: f64,
    pub upper_ram_depth: Option<f64>,
    pub tool_clearance_factor: f64, // Added to replace hardcoded 0.5 multiplier
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CollisionResponse {
    pub is_collision: bool,
    pub interference_distance: f64,
    pub recommendation: String,
}

pub struct CollisionAnalyzer;

impl CollisionAnalyzer {
    #[instrument(skip(req))]
    pub fn detect_collision(req: &CollisionRequest) -> Result<CollisionResponse, EngineError> {
        info!("Detecting bending collision for flange length: {}", req.flange_length);

        if req.flange_length <= 0.0 || req.bend_angle <= 0.0 {
            return Err(EngineError::InvalidInput("Flange length and angle must be > 0".to_string()));
        }

        let angle_rad = req.bend_angle * std::f64::consts::PI / 180.0;
        let vertical_height = req.flange_length * angle_rad.sin();
        let horizontal_encroachment = req.flange_length * (1.0 - angle_rad.cos());

        let mut collision = false;
        let mut interference = 0.0;
        let mut rec = "Safe to bend".to_string();

        let tool_clearance_width = req.v_die_width * req.tool_clearance_factor;

        if vertical_height > req.punch_height && horizontal_encroachment > tool_clearance_width {
            collision = true;
            interference = vertical_height - req.punch_height;
            rec = "Collision detected! Use a taller punch (Gooseneck punch) or bend in a different sequence.";
            warn!("{}", rec);
        }

        if let Some(ram_depth) = req.upper_ram_depth {
            if vertical_height > req.punch_height + ram_depth {
                collision = true;
                interference = vertical_height - (req.punch_height + ram_depth);
                rec = "CRITICAL: Flange hits the Upper Ram! You must change the bending sequence or use a taller tool setup.";
                warn!("{}", rec);
            }
        }

        Ok(CollisionResponse {
            is_collision: collision,
            interference_distance: interference,
            recommendation: rec.to_string(),
        })
    }
}
