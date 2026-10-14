use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct KerfRequest {
    pub material: String,
    pub thickness: f64,
    pub laser_power_watts: Option<f64>,
    pub steel_base_kerf: f64,
    pub alu_base_kerf: f64,
    pub stainless_base_kerf: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct KerfResponse {
    pub kerf_width: f64,
    pub thermal_expansion_factor: f64, // mm per meter of cutting
    pub recommended_offset: f64,
}

#[instrument(skip(req))]
pub fn calculate_kerf(req: &KerfRequest) -> Result<KerfResponse, EngineError> {
    info!("Calculating Kerf for material {} at {}mm", req.material, req.thickness);

    if req.thickness <= 0.0 {
        return Err(EngineError::InvalidInput("Thickness must be positive".to_string()));
    }

    let base_kerf = match req.material.to_lowercase().as_str() {
        "steel" | "mild steel" => req.steel_base_kerf + (req.thickness * 0.01),
        "aluminum" => req.alu_base_kerf + (req.thickness * 0.02),
        "stainless steel" => req.stainless_base_kerf + (req.thickness * 0.015),
        _ => req.steel_base_kerf,
    };

    let power_mod = if let Some(p) = req.laser_power_watts {
        if p > 4000.0 { 1.1 } else { 1.0 }
    } else {
        1.0
    };

    let actual_kerf = base_kerf * power_mod;

    let thermal_expansion = if req.thickness > 10.0 {
        0.5 
    } else {
        0.1
    };

    let recommended_offset = (actual_kerf / 2.0) - (thermal_expansion * 0.01);

    Ok(KerfResponse {
        kerf_width: actual_kerf,
        thermal_expansion_factor: thermal_expansion,
        recommended_offset,
    })
}
