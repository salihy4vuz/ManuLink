use serde::{Serialize, Deserialize};
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub enum FormType {
    Louver,
    Countersink,
    Dimple,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FormingRequest {
    pub thickness_mm: f64,
    pub material: String,
    pub form_type: FormType,
    pub quantity: u32,
    pub time_per_hit_sec: f64,
    pub hardness_factor: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FormingResult {
    pub total_time_sec: f64,
    pub recommended_tonnage_kn: f64,
    pub is_feasible: bool,
    pub warning: Option<String>,
}

pub struct FormingAnalyzer;

impl FormingAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &FormingRequest) -> FormingResult {
        info!("Analyzing forming for {:?}", req.form_type);

        let base_tonnage = 20.0 * req.thickness_mm * req.hardness_factor;
        
        let mut warning = None;
        let mut is_feasible = true;

        if req.thickness_mm > 3.0 && matches!(req.form_type, FormType::Louver) {
            is_feasible = false;
            let msg = "Louver forming on thickness > 3mm is highly likely to break the tool.".to_string();
            warn!("{}", msg);
            warning = Some(msg);
        }

        FormingResult {
            total_time_sec: req.time_per_hit_sec * req.quantity as f64,
            recommended_tonnage_kn: base_tonnage,
            is_feasible,
            warning,
        }
    }
}
