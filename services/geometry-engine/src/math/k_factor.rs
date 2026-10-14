use serde::{Deserialize, Serialize};
use crate::error::EngineError;
use tracing::{info, instrument};

#[derive(Debug, Serialize, Deserialize)]
pub struct KFactorRequest {
    pub material: String,
    pub thickness: f64,
    pub inner_radius: f64,
    pub bend_angle_degrees: f64, 
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct KFactorResponse {
    pub k_factor: f64,
    pub bend_allowance: f64,
    pub bend_deduction: f64,
    pub setback: f64,
}

pub trait KFactorCalculator {
    fn calculate(&self, req: &KFactorRequest) -> Result<KFactorResponse, EngineError>;
}

pub struct DynamicKFactorCalculator;

impl KFactorCalculator for DynamicKFactorCalculator {
    #[instrument(skip(self, req))]
    fn calculate(&self, req: &KFactorRequest) -> Result<KFactorResponse, EngineError> {
        info!("Calculating dynamic K-Factor for material {} at {}mm", req.material, req.thickness);

        if req.thickness <= 0.0 || req.inner_radius < 0.0 || req.bend_angle_degrees <= 0.0 || req.bend_angle_degrees >= 180.0 {
            return Err(EngineError::InvalidInput("Thickness and Bend Angle must be > 0. Bend angle must be < 180.".to_string()));
        }

        let rt_ratio = req.inner_radius / req.thickness;
        
        let mut k_factor: f64 = match req.material.to_lowercase().as_str() {
            "steel" | "mild steel" => 0.42,
            "aluminum" => 0.44,
            "stainless steel" | "stainless" => 0.38,
            _ => 0.40, 
        };

        if rt_ratio < 1.0 {
            k_factor -= 0.05; 
        } else if rt_ratio > 2.0 {
            k_factor += 0.05;
        }
        
        k_factor = k_factor.clamp(0.2, 0.5);

        let bend_angle_rad = req.bend_angle_degrees * (std::f64::consts::PI / 180.0);
        
        let bend_allowance = bend_angle_rad * (req.inner_radius + k_factor * req.thickness);
        
        let tan_half_angle = (bend_angle_rad / 2.0).tan();
        let setback = (req.inner_radius + req.thickness) * tan_half_angle;
        let bend_deduction = (2.0 * setback) - bend_allowance;
        
        Ok(KFactorResponse {
            k_factor,
            bend_allowance,
            bend_deduction,
            setback,
        })
    }
}
