use serde::{Serialize, Deserialize};
use tracing::{info, warn, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct ToolingRequest {
    pub material_type: String, 
    pub thickness_mm: f64,
    pub desired_bend_radius_mm: f64,
    pub part_length_meter: f64,
    pub max_machine_tonnage_kn: f64,
    
    pub tensile_factor_stainless: f64,
    pub tensile_factor_aluminum: f64,
    pub tensile_factor_mild: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ToolingRecommendation {
    pub recommended_die_v_opening_mm: f64,
    pub recommended_punch_radius_mm: f64,
    pub minimum_flange_length_mm: f64,
    pub tonnage_per_meter_kn: f64,
    pub total_required_tonnage_kn: f64,
    pub estimated_springback_angle_deg: f64,
    pub is_tonnage_exceeded: bool,
    pub warning_message: Option<String>,
}

pub struct ToolingAI;

impl ToolingAI {
    #[instrument(skip(req))]
    pub fn select_tooling(req: &ToolingRequest) -> ToolingRecommendation {
        info!("Selecting tooling for thickness {}mm, material {}", req.thickness_mm, req.material_type);

        let multiplier = if req.thickness_mm <= 3.0 {
            6.0
        } else if req.thickness_mm <= 8.0 {
            8.0
        } else {
            10.0
        };

        let mut v_opening = req.thickness_mm * multiplier;
        
        if req.material_type.to_lowercase().contains("stainless") {
            v_opening *= 1.2;
        }

        let standard_v_openings = [4.0, 6.0, 8.0, 10.0, 12.0, 16.0, 20.0, 24.0, 32.0, 40.0, 50.0];
        let mut selected_v = *standard_v_openings.last().unwrap();
        for &v in standard_v_openings.iter() {
            if v >= v_opening {
                selected_v = v;
                break;
            }
        }

        let natural_ir = selected_v / 6.0;
        let punch_radius = if req.desired_bend_radius_mm > 0.0 && req.desired_bend_radius_mm <= natural_ir {
            req.desired_bend_radius_mm
        } else {
            natural_ir.max(0.5) 
        };

        let min_flange = selected_v * 0.7;

        let tensile_factor = match req.material_type.to_lowercase().as_str() {
            "stainless steel" => req.tensile_factor_stainless,
            "aluminum" => req.tensile_factor_aluminum,
            _ => req.tensile_factor_mild, 
        };

        let tonnage_per_meter_kn = (tensile_factor * 1000.0 * req.thickness_mm.powi(2)) / selected_v / 1000.0;
        let total_tonnage_kn = tonnage_per_meter_kn * req.part_length_meter;

        let mut warning = None;
        let mut exceeded = false;
        if req.max_machine_tonnage_kn > 0.0 && total_tonnage_kn > req.max_machine_tonnage_kn {
            exceeded = true;
            let msg = format!(
                "CRITICAL: Required tonnage ({:.1} kN) exceeds machine capacity ({:.1} kN)!",
                total_tonnage_kn, req.max_machine_tonnage_kn
            );
            warn!("{}", msg);
            warning = Some(msg);
        }

        let base_springback = match req.material_type.to_lowercase().as_str() {
            "stainless steel" => 2.0,
            "aluminum" => 1.5,
            _ => 1.0, 
        };
        let springback_angle = base_springback * (punch_radius / req.thickness_mm).max(1.0);

        ToolingRecommendation {
            recommended_die_v_opening_mm: selected_v,
            recommended_punch_radius_mm: punch_radius,
            minimum_flange_length_mm: min_flange,
            tonnage_per_meter_kn,
            total_required_tonnage_kn: total_tonnage_kn,
            estimated_springback_angle_deg: springback_angle,
            is_tonnage_exceeded: exceeded,
            warning_message: warning,
        }
    }
}
