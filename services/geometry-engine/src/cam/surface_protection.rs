use serde::{Serialize, Deserialize};
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct SurfaceProtectionRequest {
    pub has_pvc_film: bool,
    pub is_brushed_or_polished: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SurfaceProtectionResult {
    pub requires_vaporization_pass: bool,
    pub recommended_die_type: String,
    pub warning: Option<String>,
}

pub struct SurfaceAnalyzer;

impl SurfaceAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &SurfaceProtectionRequest) -> SurfaceProtectionResult {
        info!("Analyzing surface protection requirements");

        let mut requires_vaporization = false;
        let mut die_type = "Standard Steel V-Die".to_string();
        let mut warning = None;

        if req.has_pvc_film {
            requires_vaporization = true;
        }

        if req.is_brushed_or_polished {
            die_type = "Urethane/Polyurethane Die".to_string();
            let msg = "CRITICAL: Material is brushed/polished. Using standard steel dies will cause bend marks (scratches). Use Urethane dies.";
            warn!("{}", msg);
            warning = Some(msg.to_string());
        }

        SurfaceProtectionResult {
            requires_vaporization_pass: requires_vaporization,
            recommended_die_type: die_type,
            warning,
        }
    }
}
