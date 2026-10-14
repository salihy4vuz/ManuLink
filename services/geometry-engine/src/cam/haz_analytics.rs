use serde::{Serialize, Deserialize};
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct HazRequest {
    pub material: String,
    pub thickness_mm: f64,
    pub has_tapped_holes: bool, 
    pub cut_process: String, 
    // Config params
    pub plasma_base_haz_mm: f64,
    pub plasma_thickness_multiplier: f64,
    pub laser_base_haz_mm: f64,
    pub laser_thickness_multiplier: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct HazResult {
    pub haz_depth_mm: f64,
    pub warning: Option<String>,
    pub recommended_drill_type: Option<String>,
}

pub struct HazAnalyzer;

impl HazAnalyzer {
    #[instrument(skip(req))]
    pub fn analyze(req: &HazRequest) -> HazResult {
        info!("Analyzing HAZ for process {}", req.cut_process);

        let base_haz = if req.cut_process.to_lowercase() == "plasma" {
            req.plasma_base_haz_mm + (req.thickness_mm * req.plasma_thickness_multiplier)
        } else {
            req.laser_base_haz_mm + (req.thickness_mm * req.laser_thickness_multiplier)
        };

        let mut warning = None;
        let mut recommended = None;

        if req.has_tapped_holes {
            if req.cut_process.to_lowercase() == "plasma" || req.material.to_lowercase().contains("stainless") {
                let w = "WARNING: Tapping inside a Laser/Plasma cut hole will break HSS taps due to edge hardening (HAZ).";
                let r = "Use Carbide (Karbür) Taps or drill holes mechanically instead of thermal cutting.";
                warn!("{}", w);
                warning = Some(w.to_string());
                recommended = Some(r.to_string());
            }
        }

        HazResult {
            haz_depth_mm: base_haz,
            warning,
            recommended_drill_type: recommended,
        }
    }
}
