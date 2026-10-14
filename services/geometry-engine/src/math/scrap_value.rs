use serde::{Serialize, Deserialize};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct ScrapValueRequest {
    pub scrap_weight_kg: f64,
    pub material_type: String,
    pub current_lme_price_eur_per_kg: f64,
    pub alu_fraction: f64,
    pub stainless_fraction: f64,
    pub mild_fraction: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ScrapValueResult {
    pub estimated_resale_value_eur: f64,
}

pub struct ScrapAnalyzer;

impl ScrapAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &ScrapValueRequest) -> ScrapValueResult {
        info!("Analyzing scrap value for material {}", req.material_type);

        let resale_fraction = match req.material_type.to_lowercase().as_str() {
            "aluminum" => req.alu_fraction,
            "stainless steel" => req.stainless_fraction,
            _ => req.mild_fraction,
        };

        let resale_value = req.scrap_weight_kg * req.current_lme_price_eur_per_kg * resale_fraction;

        ScrapValueResult {
            estimated_resale_value_eur: resale_value,
        }
    }
}
