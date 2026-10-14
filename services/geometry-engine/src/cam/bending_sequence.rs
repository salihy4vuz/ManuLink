use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, warn, error, instrument};
use std::env;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct BendOperation {
    pub id: String,
    pub angle: f64,
    pub length: f64,
    pub direction: Option<String>,
    pub tool_id: Option<String>,
}

#[derive(Debug, Clone)]
pub struct BoundingBox {
    pub x_min: f64, pub x_max: f64,
    pub y_min: f64, pub y_max: f64,
    pub z_min: f64, pub z_max: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SequenceRequest {
    pub bends: Vec<BendOperation>,
    pub max_flips_allowed: u32,
    pub punch_height_mm: f64,
    pub punch_width_mm: f64,
    pub sheet_thickness_mm: f64,
    pub sheet_width_mm: f64,
    pub sheet_length_mm: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SequenceResponse {
    pub optimal_sequence: Vec<String>,
    pub estimated_flips: u32,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct AIResponse {
    pub sequence: Vec<String>,
    pub flips: u32,
}

pub struct BendingSequenceAnalyzer;

impl BendingSequenceAnalyzer {
    #[instrument(skip(req))]
    pub async fn optimize_bending_sequence(req: &SequenceRequest) -> Result<SequenceResponse, EngineError> {
        info!("Optimizing bending sequence for {} bends", req.bends.len());

        if req.bends.is_empty() {
            return Err(EngineError::InvalidInput("At least one bend is required".to_string()));
        }

        let ai_endpoint = env::var("BENDING_AI_ENDPOINT").unwrap_or_else(|_| "http://localhost:8080/v1/models/bender:predict".to_string());
        let client = reqwest::Client::builder()
            .timeout(std::time::Duration::from_secs(3))
            .build()
            .map_err(|e| EngineError::SystemError(e.to_string()))?;

        if let Ok(res) = client.post(&ai_endpoint).json(req).send().await {
            if res.status().is_success() {
                if let Ok(ai_data) = res.json::<AIResponse>().await {
                    info!("Successfully retrieved AI-optimized sequence from {}", ai_endpoint);
                    return Ok(SequenceResponse {
                        optimal_sequence: ai_data.sequence,
                        estimated_flips: ai_data.flips,
                    });
                }
            }
        }

        warn!("[Fallback] AI Model unavailable or timed out at {}. Falling back to Heuristic Algorithm.", ai_endpoint);

        let mut sorted_bends = req.bends.clone();
        sorted_bends.sort_by(|a, b| a.length.partial_cmp(&b.length).unwrap());

        let mut sequence = Vec::new();
        let mut current_tool = "".to_string();
        let mut flips = 0;
        
        let mut state_bbox = BoundingBox {
            x_min: 0.0, x_max: req.sheet_length_mm,
            y_min: 0.0, y_max: req.sheet_width_mm,
            z_min: 0.0, z_max: req.sheet_thickness_mm,
        };

        let mut remaining = sorted_bends.clone();
        
        while !remaining.is_empty() {
            let mut best_idx = 0;
            let mut collision_detected = false;

            for (i, bend) in remaining.iter().enumerate() {
                if Self::simulate_collision(&bend, &state_bbox, req.punch_height_mm, req.punch_width_mm) {
                    warn!("Collision detected for bend {}. Re-evaluating sequence.", bend.id);
                    collision_detected = true;
                    continue; 
                }
                best_idx = i;
                break;
            }

            if collision_detected && best_idx == 0 && remaining.len() > 1 {
                warn!("Forced to proceed with bend {} despite collision risk.", remaining[0].id);
            }

            let bend = remaining.remove(best_idx);
            sequence.push(bend.id.clone());

            Self::update_state(&mut state_bbox, &bend);
            
            if let Some(t) = bend.tool_id {
                if current_tool != "" && current_tool != t {
                    flips += 1;
                }
                current_tool = t;
            } else {
                if flips < req.max_flips_allowed {
                    flips += 1;
                }
            }
        }

        info!("Heuristic sequence completed with {} flips.", flips);

        Ok(SequenceResponse {
            optimal_sequence: sequence,
            estimated_flips: flips,
        })
    }

    fn simulate_collision(bend: &BendOperation, bbox: &BoundingBox, punch_height: f64, punch_width: f64) -> bool {
        let angle_rad = bend.angle.to_radians();
        let tip_z = bend.length * angle_rad.sin();
        let tip_y = bend.length * angle_rad.cos(); 

        if tip_z > punch_height {
            return true;
        }
        
        if tip_z > 10.0 && tip_y < (punch_width / 2.0) && bend.angle > 90.0 {
            return true;
        }
        
        if bbox.z_max > punch_height {
            return true;
        }

        false
    }

    fn update_state(bbox: &mut BoundingBox, bend: &BendOperation) {
        let angle_rad = bend.angle.to_radians();
        let dir = bend.direction.as_deref().unwrap_or("UP");
        
        let delta_z = bend.length * angle_rad.sin();
        let delta_y = bend.length * (1.0 - angle_rad.cos());
        
        if dir == "UP" {
            bbox.z_max = f64::max(bbox.z_max, delta_z);
        } else {
            bbox.z_min = f64::min(bbox.z_min, -delta_z);
        }
        
        bbox.y_max -= delta_y;
    }
}
