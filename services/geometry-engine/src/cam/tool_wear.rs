use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct ToolWearRequest {
    pub process_type: String, 
    pub hits_or_pierces: u32,
    pub thickness_mm: f64,
    pub punch_wear_rate: f64,
    pub punch_cost: f64,
    pub laser_wear_rate: f64,
    pub laser_cost: f64,
    pub bending_wear_rate: f64,
    pub bending_cost: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ToolWearResult {
    pub wear_percentage: f64,
    pub estimated_maintenance_cost_eur: f64,
}

pub struct ToolWearAnalyzer;

impl ToolWearAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &ToolWearRequest) -> ToolWearResult {
        info!("Analyzing tool wear for process {}", req.process_type);

        let (wear_rate, cost_per_100_percent) = match req.process_type.to_lowercase().as_str() {
            "punch" => (req.punch_wear_rate * req.thickness_mm, req.punch_cost), 
            "laser" => (req.laser_wear_rate * req.thickness_mm, req.laser_cost), 
            "bending" => (req.bending_wear_rate * req.thickness_mm, req.bending_cost), 
            _ => (0.0001, 100.0),
        };

        let total_wear = wear_rate * req.hits_or_pierces as f64;
        let cost = total_wear * cost_per_100_percent;

        ToolWearResult {
            wear_percentage: total_wear * 100.0, 
            estimated_maintenance_cost_eur: cost,
        }
    }
}
