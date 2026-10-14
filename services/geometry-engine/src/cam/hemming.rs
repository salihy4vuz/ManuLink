use serde::{Serialize, Deserialize};
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub enum HemmingType {
    Flat, 
    TearDrop, 
    Open, 
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HemmingRequest {
    pub thickness_mm: f64,
    pub hem_length_mm: f64,
    pub part_width_mm: f64,
    pub hemming_type: HemmingType,
    // Config params
    pub flattening_force_factor: f64,
    pub base_tonnage_multiplier: f64,
    pub elongation_factor: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HemmingResult {
    pub required_tonnage_kn: f64,
    pub recommended_flattening_die: bool,
    pub elongation_allowance_mm: f64,
    pub warning: Option<String>,
}

pub struct HemmingAnalyzer;

impl HemmingAnalyzer {
    #[instrument(skip(req))]
    pub fn calculate(req: &HemmingRequest) -> HemmingResult {
        info!("Calculating hemming requirements for type {:?}", req.hemming_type);

        let base_tonnage = (req.thickness_mm * req.thickness_mm * req.base_tonnage_multiplier) / (req.thickness_mm * 2.0); 
        let required_tonnage = base_tonnage * req.flattening_force_factor * (req.part_width_mm / 1000.0);

        let mut warning = None;
        if req.hem_length_mm < req.thickness_mm * 4.0 {
            let msg = format!("Hem length ({:.1}mm) is dangerously short. Minimum recommended is {:.1}mm.", req.hem_length_mm, req.thickness_mm * 4.0);
            warn!("{}", msg);
            warning = Some(msg);
        }

        let elongation = req.thickness_mm * req.elongation_factor;

        HemmingResult {
            required_tonnage_kn: required_tonnage,
            recommended_flattening_die: true,
            elongation_allowance_mm: elongation,
            warning,
        }
    }
}
