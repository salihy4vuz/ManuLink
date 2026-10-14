use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct SustainabilityRequest {
    pub energy_kwh: f64,
    pub material_weight_kg: f64,
    pub scrap_weight_kg: f64,
    pub energy_co2_per_kwh: f64,
    pub material_co2_per_kg: f64,
    pub scrap_penalty_co2_per_kg: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SustainabilityResult {
    pub total_co2_emission_kg: f64,
    pub energy_co2_kg: f64,
    pub material_co2_kg: f64,
    pub sustainability_score: f64, 
}

pub struct SustainabilityAnalyzer;

impl SustainabilityAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &SustainabilityRequest) -> SustainabilityResult {
        info!("Analyzing sustainability metrics");

        let energy_co2 = req.energy_kwh * req.energy_co2_per_kwh;
        let material_co2 = req.material_weight_kg * req.material_co2_per_kg;
        let scrap_penalty = req.scrap_weight_kg * req.scrap_penalty_co2_per_kg;

        let total_co2 = energy_co2 + material_co2 + scrap_penalty;

        let scrap_ratio = if req.material_weight_kg > 0.0 {
            req.scrap_weight_kg / req.material_weight_kg
        } else {
            0.0
        };
        
        let mut score = 100.0 - (scrap_ratio * 200.0); 
        score = score.clamp(0.0, 100.0);

        SustainabilityResult {
            total_co2_emission_kg: total_co2,
            energy_co2_kg: energy_co2,
            material_co2_kg: material_co2,
            sustainability_score: score,
        }
    }
}
