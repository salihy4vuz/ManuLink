use serde::{Deserialize, Serialize};
use tracing::{info, instrument};
use anyhow::{Result, Context};

#[derive(Debug, Serialize, Deserialize)]
pub struct BendingRequest {
    pub material: String,
    pub thickness: f64,
    pub bend_length: f64,
    pub tensile_strength_mpa: f64, // Must be provided by DB now
    pub v_die_multiplier_small: f64,
    pub v_die_multiplier_med: f64,
    pub v_die_multiplier_large: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct BendingResponse {
    pub required_tonnage: f64,
    pub recommended_v_die: f64,
    pub minimum_flange_length: f64,
}

pub struct BendingAnalyzer;

impl BendingAnalyzer {
    #[instrument(skip(req))]
    pub fn calculate_bending_force(req: &BendingRequest) -> Result<BendingResponse> {
        info!("Calculating bending force for material: {}, thickness: {}", req.material, req.thickness);

        if req.thickness <= 0.0 || req.bend_length <= 0.0 {
            return Err(anyhow::anyhow!("Thickness and bend length must be strictly positive"));
        }

        let v_die = if req.thickness <= 3.0 {
            req.thickness * req.v_die_multiplier_small
        } else if req.thickness <= 8.0 {
            req.thickness * req.v_die_multiplier_med
        } else {
            req.thickness * req.v_die_multiplier_large
        };

        let length_in_meters = req.bend_length / 1000.0;
        let force_kn = (1.33 * length_in_meters * req.thickness.powi(2) * req.tensile_strength_mpa) / v_die;
        let tonnage = force_kn / 9.81;

        let min_flange = (v_die * 0.7).ceil();

        Ok(BendingResponse {
            required_tonnage: tonnage,
            recommended_v_die: v_die,
            minimum_flange_length: min_flange,
        })
    }
}
