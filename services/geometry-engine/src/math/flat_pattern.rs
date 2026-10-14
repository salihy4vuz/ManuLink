use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use super::k_factor::{KFactorRequest, DynamicKFactorCalculator, KFactorCalculator};
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct FlatPatternRequest {
    pub material: String,
    pub thickness: f64,
    pub flanges: Vec<f64>, // Lengths of straight flanges
    pub bend_angles: Vec<f64>, // Angles in degrees for each bend
    pub inner_radius: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct FlatPatternResponse {
    pub total_flat_length: f64,
    pub bend_deductions_used: Vec<f64>,
}

#[instrument(skip(req))]
pub fn calculate_flat_pattern(req: &FlatPatternRequest) -> Result<FlatPatternResponse, EngineError> {
    info!("Calculating flat pattern for {} bends", req.bend_angles.len());
    if req.flanges.len() != req.bend_angles.len() + 1 {
        return Err(EngineError::InvalidInput(
            "Number of flanges must be exactly one more than the number of bends.".to_string()
        ));
    }

    let k_calc = DynamicKFactorCalculator;
    let mut total_length: f64 = req.flanges.iter().sum();
    let mut deductions = Vec::new();

    // For each bend, calculate specific Bend Deduction based on its angle and subtract it
    for angle in &req.bend_angles {
        let k_req = KFactorRequest {
            material: req.material.clone(),
            thickness: req.thickness,
            inner_radius: req.inner_radius,
            bend_angle_degrees: *angle,
        };
        
        let k_resp = k_calc.calculate(&k_req)?;
        total_length -= k_resp.bend_deduction;
        deductions.push(k_resp.bend_deduction);
    }
    
    Ok(FlatPatternResponse {
        total_flat_length: total_length,
        bend_deductions_used: deductions,
    })
}
